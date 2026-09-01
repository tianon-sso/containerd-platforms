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

import "golang.org/x/sys/unix"

// rva22Bits and rva23Bits are the subset of each RVA profile's mandatory
// extensions —
// https://github.com/riscv/riscv-isa-manual/blob/07531fdab1b9e2c3c2306fd3805dfeea1a39b42f/src/profiles/rva22.adoc,
// https://github.com/riscv/riscv-isa-manual/blob/07531fdab1b9e2c3c2306fd3805dfeea1a39b42f/src/profiles/rva23.adoc —
// that riscv_hwprobe can actually report. Each profile also mandates a
// number of memory-model/cache-coherency guarantees (e.g. Ziccif, Za64rs,
// Zic64b) that have no corresponding hwprobe bit at all — those are assumed
// to hold for any Linux-capable riscv64 hart, the same way x86-64-v1's
// baseline (cmov, fpu, mmx, ...) is never itself probed.
const (
	rva22Bits = unix.RISCV_HWPROBE_EXT_ZBA |
		unix.RISCV_HWPROBE_EXT_ZBB |
		unix.RISCV_HWPROBE_EXT_ZBS |
		unix.RISCV_HWPROBE_EXT_ZIHINTPAUSE

	rva23Bits = unix.RISCV_HWPROBE_IMA_V |
		unix.RISCV_HWPROBE_EXT_ZVFHMIN |
		unix.RISCV_HWPROBE_EXT_ZVBB |
		unix.RISCV_HWPROBE_EXT_ZVKT |
		unix.RISCV_HWPROBE_EXT_ZIHINTNTL |
		unix.RISCV_HWPROBE_EXT_ZICOND
)

// getRISCV64Variant returns the highest RVA profile ("rva20u64", "rva22u64",
// "rva23u64" — see https://github.com/riscv/riscv-isa-manual/tree/07531fdab1b9e2c3c2306fd3805dfeea1a39b42f/src/profiles)
// the host CPU satisfies, or "" if detection found nothing conclusive (e.g.
// a pre-6.4 kernel, which predates the riscv_hwprobe(2) syscall this uses).
// Detection is a syscall via golang.org/x/sys/unix; no assembly is used.
func getRISCV64Variant() (string, error) {
	pairs := []unix.RISCVHWProbePairs{
		{Key: unix.RISCV_HWPROBE_KEY_IMA_EXT_0},
	}
	if err := unix.RISCVHWProbe(pairs, nil, 0); err != nil {
		return "", err
	}
	if pairs[0].Key == -1 {
		// Kernel doesn't recognize this key: no reliable signal.
		return "", nil
	}
	return riscvProfileFromExtBits(pairs[0].Value), nil
}

// riscvProfileFromExtBits maps a RISCV_HWPROBE_KEY_IMA_EXT_0 bitmask to the
// highest RVA profile it satisfies (see rva22Bits/rva23Bits above).
func riscvProfileFromExtBits(ext uint64) string {
	switch {
	case ext&rva23Bits == rva23Bits:
		return "rva23u64"
	case ext&rva22Bits == rva22Bits:
		return "rva22u64"
	}
	return "rva20u64"
}
