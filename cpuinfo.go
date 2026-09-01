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
	"runtime"
	"sync"

	"github.com/containerd/log"
)

// cpuVariant returns the detected CPU variant, e.g. v5/v6/v7/v8 (arm),
// v2/v3/v4 (amd64), power8/power9/power10 (ppc64le),
// rva20u64/rva22u64/rva23u64 (riscv64). Empty if the architecture has no
// variant concept, or detection found nothing conclusive.
var cpuVariant = sync.OnceValue(func() string {
	var (
		variant string
		err     error
	)
	switch runtime.GOARCH {
	case "arm", "arm64":
		variant, err = getARMVariant()
	case "amd64":
		variant, err = getAMD64Variant()
	case "ppc64le":
		variant, err = getPPC64LEVariant()
	case "riscv64":
		variant, err = getRISCV64Variant()
	default:
		return ""
	}
	if err != nil {
		log.L.Errorf("Error detecting CPU variant for %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
	}
	return variant
})
