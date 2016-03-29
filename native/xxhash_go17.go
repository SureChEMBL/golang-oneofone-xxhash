// +build go1.7 appengine safe

package xxhash

//const Backend = "GoSafe17"

// Checksum32S returns the checksum of the input bytes with the specific seed.
func Checksum32S(in []byte, seed uint32) (h uint32) {
	i, l := 0, len(in)
	if l >= 16 {
		var (
			v1 = seed + prime32x1 + prime32x2
			v2 = seed + prime32x2
			v3 = seed + 0
			v4 = seed - prime32x1
		)
		for ; i <= l-16; i += 16 {
			in := in[i : i+16 : l]
			v1 += (uint32(in[0]) | uint32(in[1])<<8 | uint32(in[2])<<16 | uint32(in[3])<<24) * prime32x2
			v1 = rotl32_13(v1) * prime32x1

			v2 += (uint32(in[4]) | uint32(in[5])<<8 | uint32(in[6])<<16 | uint32(in[7])<<24) * prime32x2
			v2 = rotl32_13(v2) * prime32x1

			v3 += (uint32(in[8]) | uint32(in[9])<<8 | uint32(in[10])<<16 | uint32(in[11])<<24) * prime32x2
			v3 = rotl32_13(v3) * prime32x1

			v4 += (uint32(in[12]) | uint32(in[13])<<8 | uint32(in[14])<<16 | uint32(in[15])<<24) * prime32x2
			v4 = rotl32_13(v4) * prime32x1
		}

		h = rotl32_1(v1) + rotl32_7(v2) + rotl32_12(v3) + rotl32_18(v4)

	} else {
		h = seed + prime32x5
	}

	h += uint32(l)
	for ; i <= l-4; i += 4 {
		in := in[i : i+4 : l]
		h += (uint32(in[0]) | uint32(in[1])<<8 | uint32(in[2])<<16 | uint32(in[3])<<24) * prime32x3
		h = rotl32_17(h) * prime32x4
	}

	for ; i < l; i++ {
		h += uint32(in[i]) * prime32x5
		h = rotl32_11(h) * prime32x1
	}

	h ^= h >> 15
	h *= prime32x2
	h ^= h >> 13
	h *= prime32x3
	h ^= h >> 16

	return
}

func (xx *XXHash32) Write(in []byte) (n int, err error) {
	i, l, ml := int32(0), int32(len(in)), xx.memIdx
	xx.ln += int32(l)

	if d := 16 - ml; ml > 0 && ml+l > 16 {
		xx.memIdx += int32(copy(xx.mem[xx.memIdx:], in[:d]))
		in = in[d:]
		ml, l = 16, int32(len(in))
	} else if ml+l < 16 {
		xx.memIdx += int32(copy(xx.mem[xx.memIdx:], in))
		return int(l), nil
	}

	if ml > 0 {
		i += 16 - ml
		xx.memIdx += int32(copy(xx.mem[xx.memIdx:], in))
		in := xx.mem[:16]

		xx.v1 += (uint32(in[0]) | uint32(in[1])<<8 | uint32(in[2])<<16 | uint32(in[3])<<24) * prime32x2
		xx.v1 = rotl32_13(xx.v1) * prime32x1

		xx.v2 += (uint32(in[4]) | uint32(in[5])<<8 | uint32(in[6])<<16 | uint32(in[7])<<24) * prime32x2
		xx.v2 = rotl32_13(xx.v2) * prime32x1

		xx.v3 += (uint32(in[8]) | uint32(in[9])<<8 | uint32(in[10])<<16 | uint32(in[11])<<24) * prime32x2
		xx.v3 = rotl32_13(xx.v3) * prime32x1

		xx.v4 += (uint32(in[12]) | uint32(in[13])<<8 | uint32(in[14])<<16 | uint32(in[15])<<24) * prime32x2
		xx.v4 = rotl32_13(xx.v4) * prime32x1

		xx.memIdx = 0
	}

	if l >= 16 {
		for ; i <= l-16; i += 16 {
			in := in[i : i+16 : l]
			xx.v1 += (uint32(in[0]) | uint32(in[1])<<8 | uint32(in[2])<<16 | uint32(in[3])<<24) * prime32x2
			xx.v1 = rotl32_13(xx.v1) * prime32x1

			xx.v2 += (uint32(in[4]) | uint32(in[5])<<8 | uint32(in[6])<<16 | uint32(in[7])<<24) * prime32x2
			xx.v2 = rotl32_13(xx.v2) * prime32x1

			xx.v3 += (uint32(in[8]) | uint32(in[9])<<8 | uint32(in[10])<<16 | uint32(in[11])<<24) * prime32x2
			xx.v3 = rotl32_13(xx.v3) * prime32x1

			xx.v4 += (uint32(in[12]) | uint32(in[13])<<8 | uint32(in[14])<<16 | uint32(in[15])<<24) * prime32x2
			xx.v4 = rotl32_13(xx.v4) * prime32x1
		}
	}

	if l-i != 0 {
		xx.memIdx += int32(copy(xx.mem[xx.memIdx:], in[i:]))
	}

	if debug {
		if l-i > 16 {
			panic("len(in) - i > 16")
		}

		if cap(xx.mem) > 16 {
			panic("cap(xx.mem) > 16")
		}
	}

	return int(l), nil
}

func (xx *XXHash32) Sum32() (h uint32) {
	i, l := int32(0), xx.memIdx
	if xx.ln >= 16 {
		h = rotl32_1(xx.v1) + rotl32_7(xx.v2) + rotl32_12(xx.v3) + rotl32_18(xx.v4)
	} else {
		h = xx.seed + prime32x5
	}

	h += uint32(xx.ln)

	if xx.memIdx > 0 {
		for ; i <= xx.memIdx-4; i += 4 {
			in := xx.mem[i : i+4 : l]
			h += (uint32(in[0]) | uint32(in[1])<<8 | uint32(in[2])<<16 | uint32(in[3])<<24) * prime32x3
			h = rotl32_17(h) * prime32x4
		}

		for ; i < xx.memIdx; i++ {
			h += uint32(xx.mem[i]) * prime32x5
			h = rotl32_11(h) * prime32x1
		}
	}
	h ^= h >> 15
	h *= prime32x2
	h ^= h >> 13
	h *= prime32x3
	h ^= h >> 16

	return
}

// Checksum64S returns the 64bit xxhash checksum for a single input
func Checksum64S(in []byte, seed uint64) (h uint64) {
	i, l := 0, len(in)
	if l >= 32 {
		var (
			v1 = seed + prime64x1 + prime64x2
			v2 = seed + prime64x2
			v3 = seed + 0
			v4 = seed - prime64x1
			l  = len(in)
		)
		for ; i <= l-32; i += 32 {
			in := in[i : i+32 : l]
			v1 += (uint64(in[0]) | uint64(in[1])<<8 | uint64(in[2])<<16 | uint64(in[3])<<24 |
				uint64(in[4])<<32 | uint64(in[5])<<40 | uint64(in[6])<<48 | uint64(in[7])<<56) * prime64x2
			v1 = rotl64_31(v1) * prime64x1

			v2 += (uint64(in[8]) | uint64(in[9])<<8 | uint64(in[10])<<16 | uint64(in[11])<<24 |
				uint64(in[12])<<32 | uint64(in[13])<<40 | uint64(in[14])<<48 | uint64(in[15])<<56) * prime64x2
			v2 = rotl64_31(v2) * prime64x1

			v3 += (uint64(in[16]) | uint64(in[17])<<8 | uint64(in[18])<<16 | uint64(in[19])<<24 |
				uint64(in[20])<<32 | uint64(in[21])<<40 | uint64(in[22])<<48 | uint64(in[23])<<56) * prime64x2
			v3 = rotl64_31(v3) * prime64x1

			v4 += (uint64(in[24]) | uint64(in[25])<<8 | uint64(in[26])<<16 | uint64(in[27])<<24 |
				uint64(in[28])<<32 | uint64(in[29])<<40 | uint64(in[30])<<48 | uint64(in[31])<<56) * prime64x2
			v4 = rotl64_31(v4) * prime64x1
		}

		h = rotl64_1(v1) + rotl64_7(v2) + rotl64_12(v3) + rotl64_18(v4)
		v1 *= prime64x2
		v1 = rotl64_31(v1)
		v1 *= prime64x1
		h ^= v1
		h = h*prime64x1 + prime64x4

		v2 *= prime64x2
		v2 = rotl64_31(v2)
		v2 *= prime64x1
		h ^= v2
		h = h*prime64x1 + prime64x4

		v3 *= prime64x2
		v3 = rotl64_31(v3)
		v3 *= prime64x1
		h ^= v3
		h = h*prime64x1 + prime64x4

		v4 *= prime64x2
		v4 = rotl64_31(v4)
		v4 *= prime64x1
		h ^= v4
		h = h*prime64x1 + prime64x4
	} else {
		h = seed + prime64x5
	}

	h += uint64(l)

	for ; i <= l-8; i += 8 {
		in := in[i : i+8 : l]
		k := (uint64(in[0]) | uint64(in[1])<<8 | uint64(in[2])<<16 | uint64(in[3])<<24 |
			uint64(in[4])<<32 | uint64(in[5])<<40 | uint64(in[6])<<48 | uint64(in[7])<<56)
		k *= prime64x2
		k = rotl64_31(k)
		k *= prime64x1
		h ^= k
		h = rotl64_27(h)*prime64x1 + prime64x4
	}

	for ; i <= l-4; i += 4 {
		in := in[i : i+4 : l]
		h ^= (uint64(in[0]) | uint64(in[1])<<8 | uint64(in[2])<<16 | uint64(in[3])<<24) * prime64x1
		h = rotl64_23(h)*prime64x2 + prime64x3
	}

	for ; i < l; i++ {
		h ^= uint64(in[i]) * prime64x5
		h = rotl64_11(h) * prime64x1
	}

	h ^= h >> 33
	h *= prime64x2
	h ^= h >> 29
	h *= prime64x3
	h ^= h >> 32

	return h
}

func (xx *XXHash64) Write(in []byte) (n int, err error) {
	var i, l, ml int32 = int32(0), int32(len(in)), xx.memIdx
	xx.ln += int32(l)
	if d := 32 - ml; ml > 0 && ml+l > 32 {
		xx.memIdx += int32(copy(xx.mem[xx.memIdx:], in[:d]))
		in = in[d:]
		ml, l = 32, int32(len(in))
	} else if ml+l < 32 {
		xx.memIdx += int32(copy(xx.mem[xx.memIdx:], in))
		return int(l), nil
	}

	if ml > 0 {
		i += 32 - ml
		xx.memIdx += int32(copy(xx.mem[xx.memIdx:], in))
		in := xx.mem[:]

		xx.v1 += (uint64(in[0]) | uint64(in[1])<<8 | uint64(in[2])<<16 | uint64(in[3])<<24 |
			uint64(in[4])<<32 | uint64(in[5])<<40 | uint64(in[6])<<48 | uint64(in[7])<<56) * prime64x2
		xx.v1 = rotl64_31(xx.v1) * prime64x1

		xx.v2 += (uint64(in[8]) | uint64(in[9])<<8 | uint64(in[10])<<16 | uint64(in[11])<<24 |
			uint64(in[12])<<32 | uint64(in[13])<<40 | uint64(in[14])<<48 | uint64(in[15])<<56) * prime64x2
		xx.v2 = rotl64_31(xx.v2) * prime64x1

		xx.v3 += (uint64(in[16]) | uint64(in[17])<<8 | uint64(in[18])<<16 | uint64(in[19])<<24 |
			uint64(in[20])<<32 | uint64(in[21])<<40 | uint64(in[22])<<48 | uint64(in[23])<<56) * prime64x2
		xx.v3 = rotl64_31(xx.v3) * prime64x1

		xx.v4 += (uint64(in[24]) | uint64(in[25])<<8 | uint64(in[26])<<16 | uint64(in[27])<<24 |
			uint64(in[28])<<32 | uint64(in[29])<<40 | uint64(in[30])<<48 | uint64(in[31])<<56) * prime64x2
		xx.v4 = rotl64_31(xx.v4) * prime64x1

		xx.memIdx = 0
	}

	if l >= 32 {
		for ; i <= l-32; i += 32 {
			in := in[i : i+32 : l]
			xx.v1 += (uint64(in[0]) | uint64(in[1])<<8 | uint64(in[2])<<16 | uint64(in[3])<<24 |
				uint64(in[4])<<32 | uint64(in[5])<<40 | uint64(in[6])<<48 | uint64(in[7])<<56) * prime64x2
			xx.v1 = rotl64_31(xx.v1) * prime64x1

			xx.v2 += (uint64(in[8]) | uint64(in[9])<<8 | uint64(in[10])<<16 | uint64(in[11])<<24 |
				uint64(in[12])<<32 | uint64(in[13])<<40 | uint64(in[14])<<48 | uint64(in[15])<<56) * prime64x2
			xx.v2 = rotl64_31(xx.v2) * prime64x1

			xx.v3 += (uint64(in[16]) | uint64(in[17])<<8 | uint64(in[18])<<16 | uint64(in[19])<<24 |
				uint64(in[20])<<32 | uint64(in[21])<<40 | uint64(in[22])<<48 | uint64(in[23])<<56) * prime64x2
			xx.v3 = rotl64_31(xx.v3) * prime64x1

			xx.v4 += (uint64(in[24]) | uint64(in[25])<<8 | uint64(in[26])<<16 | uint64(in[27])<<24 |
				uint64(in[28])<<32 | uint64(in[29])<<40 | uint64(in[30])<<48 | uint64(in[31])<<56) * prime64x2
			xx.v4 = rotl64_31(xx.v4) * prime64x1
		}

	}

	if l-i != 0 {
		xx.memIdx += int32(copy(xx.mem[xx.memIdx:], in[i:]))
	}

	if debug {
		if l-i > 32 {
			panic("len(in) - i > 32")
		}
	}
	return int(l), nil
}

func (xx *XXHash64) Sum64() (h uint64) {
	var i int32
	v1, v2, v3, v4 := xx.v1, xx.v2, xx.v3, xx.v4
	if xx.ln >= 32 {
		h = rotl64_1(v1) + rotl64_7(v2) + rotl64_12(v3) + rotl64_18(v4)

		v1 *= prime64x2
		v1 = rotl64_31(v1)
		v1 *= prime64x1
		h ^= v1
		h = h*prime64x1 + prime64x4

		v2 *= prime64x2
		v2 = rotl64_31(v2)
		v2 *= prime64x1
		h ^= v2
		h = h*prime64x1 + prime64x4

		v3 *= prime64x2
		v3 = rotl64_31(v3)
		v3 *= prime64x1
		h ^= v3
		h = h*prime64x1 + prime64x4

		v4 *= prime64x2
		v4 = rotl64_31(v4)
		v4 *= prime64x1
		h ^= v4
		h = h*prime64x1 + prime64x4
	} else {
		h = xx.seed + prime64x5
	}

	h += uint64(xx.ln)
	if xx.memIdx > 0 {
		in := xx.mem[:xx.memIdx]
		for ; i <= xx.memIdx-8; i += 8 {
			in := in[i : i+8 : 32]
			k := (uint64(in[0]) | uint64(in[1])<<8 | uint64(in[2])<<16 | uint64(in[3])<<24 |
				uint64(in[4])<<32 | uint64(in[5])<<40 | uint64(in[6])<<48 | uint64(in[7])<<56)
			k *= prime64x2
			k = rotl64_31(k)
			k *= prime64x1
			h ^= k
			h = rotl64_27(h)*prime64x1 + prime64x4
		}

		for ; i <= xx.memIdx-4; i += 4 {
			in := in[i : i+4 : 32]
			h ^= (uint64(in[0]) | uint64(in[1])<<8 | uint64(in[2])<<16 | uint64(in[3])<<24) * prime64x1
			h = rotl64_23(h)*prime64x2 + prime64x3
		}

		for ; i < xx.memIdx; i++ {
			h ^= uint64(in[i]) * prime64x5
			h = rotl64_11(h) * prime64x1
		}
		xx.memIdx = 0
	}
	h ^= h >> 33
	h *= prime64x2
	h ^= h >> 29
	h *= prime64x3
	h ^= h >> 32

	return
}