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
// name is there specifically to keep Go from inferring one from the "amd64"
// preceding it: cpu.X86 is a safe, always-zero-valued struct on every
// platform (see golang.org/x/sys/cpu), and cpuinfo.go calls getAMD64Variant
// unconditionally regardless of runtime.GOARCH, so this must compile
// everywhere.
import "golang.org/x/sys/cpu"

// getAMD64Variant returns the highest x86-64 microarchitecture level ("v2",
// "v3", "v4") the host CPU satisfies, per the x86-64 psABI
// (https://gitlab.com/x86-psABIs/x86-64-ABI/-/blob/e1ce098331da5dbd66e1ffc74162380bcc213236/x86-64-ABI/low-level-sys-info.tex).
// See also https://github.com/golang/go/issues/58015, which proposes
// exposing this same level directly from golang.org/x/sys/cpu.
// Returns "" if the host only satisfies the v1 baseline. Feature bits come
// from golang.org/x/sys/cpu, which reads them via the CPUID instruction; no
// assembly is implemented in this package.
func getAMD64Variant() (string, error) {
	var (
		v2 = cpu.X86.HasCX16 &&
			cpu.X86.HasLAHF &&
			cpu.X86.HasPOPCNT &&
			cpu.X86.HasSSE41 &&
			cpu.X86.HasSSE42 &&
			cpu.X86.HasSSSE3

		v3 = v2 &&
			cpu.X86.HasAVX &&
			cpu.X86.HasAVX2 &&
			cpu.X86.HasBMI1 &&
			cpu.X86.HasBMI2 &&
			cpu.X86.HasF16C &&
			cpu.X86.HasFMA &&
			cpu.X86.HasLZCNT &&
			cpu.X86.HasMOVBE &&
			cpu.X86.HasOSXSAVE &&
			cpu.X86.HasXSAVE

		v4 = v3 &&
			cpu.X86.HasAVX512F &&
			cpu.X86.HasAVX512BW &&
			cpu.X86.HasAVX512CD &&
			cpu.X86.HasAVX512DQ &&
			cpu.X86.HasAVX512VL
	)
	switch {
	case v4:
		return "v4", nil
	case v3:
		return "v3", nil
	case v2:
		return "v2", nil
	}
	return "", nil
}
