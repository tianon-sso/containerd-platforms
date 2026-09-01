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

func withPPC64Flags(t *testing.T, isPOWER8, isPOWER9, isPOWER10 bool) {
	t.Helper()
	orig := cpu.PPC64
	t.Cleanup(func() { cpu.PPC64 = orig })

	cpu.PPC64.IsPOWER8 = isPOWER8
	cpu.PPC64.IsPOWER9 = isPOWER9
	cpu.PPC64.IsPOWER10 = isPOWER10
}

func TestGetPPC64LEVariant(t *testing.T) {
	for _, tc := range []struct {
		name                          string
		isPOWER8, isPOWER9, isPOWER10 bool
		want                          string
	}{
		{name: "nothing detected", want: ""},
		{name: "power8", isPOWER8: true, want: "power8"},
		// The kernel sets every applicable lower-level HWCAP2 bit
		// cumulatively, so a real POWER9 host reports both bits set.
		{name: "power9 (implies power8)", isPOWER8: true, isPOWER9: true, want: "power9"},
		{name: "power10/power11 (implies power9, power8)", isPOWER8: true, isPOWER9: true, isPOWER10: true, want: "power10"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			withPPC64Flags(t, tc.isPOWER8, tc.isPOWER9, tc.isPOWER10)
			got, err := getPPC64LEVariant()
			if err != nil {
				t.Fatalf("getPPC64LEVariant() error = %v", err)
			}
			if got != tc.want {
				t.Errorf("getPPC64LEVariant() = %q, want %q", got, tc.want)
			}
		})
	}
}
