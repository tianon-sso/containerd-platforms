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

// This file intentionally has no GOARCH build tag, and the "_variant" in its
// name is there specifically to keep Go from inferring one from the
// "ppc64le" preceding it: cpu.PPC64 is a safe, always-zero-valued struct on
// every platform (see golang.org/x/sys/cpu), and cpuinfo.go calls
// getPPC64LEVariant unconditionally regardless of runtime.GOARCH, so this
// must compile everywhere.
import "golang.org/x/sys/cpu"

// getPPC64LEVariant returns the highest POWER ISA level ("power8", "power9",
// "power10") the host CPU satisfies, or "" if detection found nothing
// conclusive (e.g. not running on Linux, where these HWCAP2 bits come from:
// https://github.com/torvalds/linux/blob/v6.12/arch/powerpc/include/uapi/asm/cputable.h).
// POWER11 shares POWER10's ISA level (3.1) and has no separate architected
// HWCAP bit of its own — confirmed by the kernel commit that added POWER11
// support (https://github.com/torvalds/linux/commit/c2ed087ed35ca569d8179924ba560be248c758e5:
// "CPU, MMU and user (ELF_HWCAP) features are unchanged vs P10") — so a
// POWER11 host is correctly (not approximately) reported as "power10" —
// that's what it actually is.
func getPPC64LEVariant() (string, error) {
	switch {
	case cpu.PPC64.IsPOWER10:
		return "power10", nil
	case cpu.PPC64.IsPOWER9:
		return "power9", nil
	case cpu.PPC64.IsPOWER8:
		return "power8", nil
	}
	return "", nil
}
