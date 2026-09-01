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

import (
	"testing"

	"golang.org/x/sys/cpu"
)

// amd64Flags lists exactly the cpu.X86 fields getAMD64Variant reads, so a
// test can set every one of them explicitly instead of needing to "zero"
// cpu.X86 (an anonymous struct type with no exported name to construct a
// literal of).
type amd64Flags struct {
	cx16, lahf, popcnt, sse41, sse42, ssse3                        bool
	avx, avx2, bmi1, bmi2, f16c, fma, lzcnt, movbe, osxsave, xsave bool
	avx512f, avx512bw, avx512cd, avx512dq, avx512vl                bool
}

// withAMD64Flags applies f to cpu.X86 for the duration of the test,
// restoring the original (real, detected) value afterward.
func withAMD64Flags(t *testing.T, f amd64Flags) {
	t.Helper()
	// cpu.X86 is process-global; borrow t.Setenv's documented "cannot be
	// used in parallel tests" guard so misuse with t.Parallel() panics.
	t.Setenv("_amd64FlagsGuard", "")

	orig := cpu.X86
	t.Cleanup(func() { cpu.X86 = orig })

	cpu.X86.HasCX16 = f.cx16
	cpu.X86.HasLAHF = f.lahf
	cpu.X86.HasPOPCNT = f.popcnt
	cpu.X86.HasSSE41 = f.sse41
	cpu.X86.HasSSE42 = f.sse42
	cpu.X86.HasSSSE3 = f.ssse3
	cpu.X86.HasAVX = f.avx
	cpu.X86.HasAVX2 = f.avx2
	cpu.X86.HasBMI1 = f.bmi1
	cpu.X86.HasBMI2 = f.bmi2
	cpu.X86.HasF16C = f.f16c
	cpu.X86.HasFMA = f.fma
	cpu.X86.HasLZCNT = f.lzcnt
	cpu.X86.HasMOVBE = f.movbe
	cpu.X86.HasOSXSAVE = f.osxsave
	cpu.X86.HasXSAVE = f.xsave
	cpu.X86.HasAVX512F = f.avx512f
	cpu.X86.HasAVX512BW = f.avx512bw
	cpu.X86.HasAVX512CD = f.avx512cd
	cpu.X86.HasAVX512DQ = f.avx512dq
	cpu.X86.HasAVX512VL = f.avx512vl
}

var fullV2 = amd64Flags{cx16: true, lahf: true, popcnt: true, sse41: true, sse42: true, ssse3: true}

var fullV3 = func() amd64Flags {
	f := fullV2
	f.avx, f.avx2, f.bmi1, f.bmi2, f.f16c, f.fma, f.lzcnt, f.movbe, f.osxsave, f.xsave = true, true, true, true, true, true, true, true, true, true
	return f
}()

var fullV4 = func() amd64Flags {
	f := fullV3
	f.avx512f, f.avx512bw, f.avx512cd, f.avx512dq, f.avx512vl = true, true, true, true, true
	return f
}()

func TestGetAMD64Variant(t *testing.T) {
	for _, tc := range []struct {
		name  string
		flags amd64Flags
		want  string
	}{
		{name: "v1 baseline (nothing set)", flags: amd64Flags{}, want: ""},
		{name: "full v2", flags: fullV2, want: "v2"},
		{name: "full v3", flags: fullV3, want: "v3"},
		{name: "full v4", flags: fullV4, want: "v4"},
		{
			name: "v3 missing one flag (LZCNT) falls back to v2",
			flags: func() amd64Flags {
				f := fullV3
				f.lzcnt = false
				return f
			}(),
			want: "v2",
		},
		{
			name: "v4 missing one flag (AVX512VL) falls back to v3",
			flags: func() amd64Flags {
				f := fullV4
				f.avx512vl = false
				return f
			}(),
			want: "v3",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			withAMD64Flags(t, tc.flags)
			got, err := getAMD64Variant()
			if err != nil {
				t.Fatalf("getAMD64Variant() error = %v", err)
			}
			if got != tc.want {
				t.Errorf("getAMD64Variant() = %q, want %q", got, tc.want)
			}
		})
	}
}
