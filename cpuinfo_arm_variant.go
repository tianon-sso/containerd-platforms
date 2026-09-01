//go:build linux && (arm || arm64)

/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package platforms

// This file intentionally does not rely on filename-based build tag
// inference: a "cpuinfo_linux_arm.go" name would make Go infer linux&&arm
// only, silently dropping arm64 from the build. The "_variant" suffix keeps
// the filename from matching any GOARCH, so the explicit tag above is the
// only constraint in effect; cpuinfo_arm_other.go carries its exact negation.

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/sys/unix"
)

// getMachineArch retrieves the machine architecture through system call
func getMachineArch() (string, error) {
	var uname unix.Utsname
	err := unix.Uname(&uname)
	if err != nil {
		return "", err
	}

	arch := string(uname.Machine[:bytes.IndexByte(uname.Machine[:], 0)])

	return arch, nil
}

// For Linux, the kernel has already detected the ABI, ISA and Features.
// So we don't need to access the ARM registers to detect platform information
// by ourselves. We can just parse these information from /proc/cpuinfo
func getCPUInfo(pattern string) (info string, err error) {
	cpuinfo, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer cpuinfo.Close()

	// Start to Parse the Cpuinfo line by line. For SMP SoC, we parse
	// the first core is enough.
	scanner := bufio.NewScanner(cpuinfo)
	for scanner.Scan() {
		newline := scanner.Text()
		list := strings.Split(newline, ":")

		if len(list) > 1 && strings.EqualFold(strings.TrimSpace(list[0]), pattern) {
			return strings.TrimSpace(list[1]), nil
		}
	}

	// Check whether the scanner encountered errors
	err = scanner.Err()
	if err != nil {
		return "", err
	}

	return "", fmt.Errorf("getCPUInfo for pattern %s: %w", pattern, errNotFound)
}

// getARMVariantFromArch gets CPU variant from the uname machine field
// (e.g. "armv7l", "aarch64"), used when /proc/cpuinfo has no "Cpu
// architecture" field.
func getARMVariantFromArch(arch string) (string, error) {
	arch = strings.ToLower(arch)

	if arch == "aarch64" {
		return "v8", nil
	}
	if arch[0:4] == "armv" && len(arch) >= 5 {
		// Valid arch format is in form of armvXx
		switch arch[3:5] {
		case "v8":
			return "v8", nil
		case "v7":
			return "v7", nil
		case "v6":
			return "v6", nil
		case "v5":
			return "v5", nil
		default:
			return "unknown", nil
		}
	}
	return "", fmt.Errorf("getARMVariantFromArch invalid arch: %s, %w", arch, errInvalidArgument)
}

// normalizeCPUArchitecture maps /proc/cpuinfo's "Cpu architecture" field to
// a "vN" variant. Only for values read from that field; getARMVariantFromArch
// normalizes the uname fallback separately.
func normalizeCPUArchitecture(variant string) string {
	switch strings.ToLower(variant) {
	case "8", "aarch64":
		return "v8"
	case "7", "7m", "?(12)", "?(13)", "?(14)", "?(15)", "?(16)", "?(17)":
		return "v7"
	case "6", "6tej":
		return "v6"
	case "5", "5t", "5te", "5tej":
		return "v5"
	default:
		return "unknown"
	}
}

// getARMVariant returns cpu variant for ARM
// We first try reading "Cpu architecture" field from /proc/cpuinfo
// If we can't find it, then fall back using a system call
// This is to cover running ARM in emulated environment on x86 host as this field in /proc/cpuinfo
// was not present.
func getARMVariant() (string, error) {
	variant, err := getCPUInfo("Cpu architecture")
	if err != nil {
		if errors.Is(err, errNotFound) {
			// Let's try getting CPU variant from machine architecture
			arch, err := getMachineArch()
			if err != nil {
				return "", fmt.Errorf("failure getting machine architecture: %v", err)
			}

			return getARMVariantFromArch(arch)
		}
		return "", fmt.Errorf("failure getting CPU variant: %v", err)
	}

	// handle edge case for Raspberry Pi ARMv6 devices (which due to a kernel quirk, report "CPU architecture: 7")
	// https://www.raspberrypi.org/forums/viewtopic.php?t=12614
	if runtime.GOARCH == "arm" && variant == "7" {
		model, err := getCPUInfo("model name")
		if err == nil && strings.HasPrefix(strings.ToLower(model), "armv6-compatible") {
			variant = "6"
		}
	}

	return normalizeCPUArchitecture(variant), nil
}
