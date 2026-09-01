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

	"golang.org/x/sys/unix"
)

func TestRiscvProfileFromExtBits(t *testing.T) {
	for _, tc := range []struct {
		name string
		ext  uint64
		want string
	}{
		{name: "nothing set", ext: 0, want: "rva20u64"},
		{name: "full rva22 bits", ext: rva22Bits, want: "rva22u64"},
		{name: "full rva23 bits (implies rva22)", ext: rva23Bits | rva22Bits, want: "rva23u64"},
		{
			name: "rva22 missing one bit (Zihintpause) falls back to rva20",
			ext:  rva22Bits &^ unix.RISCV_HWPROBE_EXT_ZIHINTPAUSE,
			want: "rva20u64",
		},
		{
			name: "rva23 missing one bit (Zicond) falls back to rva22",
			ext:  (rva23Bits &^ unix.RISCV_HWPROBE_EXT_ZICOND) | rva22Bits,
			want: "rva22u64",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := riscvProfileFromExtBits(tc.ext); got != tc.want {
				t.Errorf("riscvProfileFromExtBits(%#x) = %q, want %q", tc.ext, got, tc.want)
			}
		})
	}
}
