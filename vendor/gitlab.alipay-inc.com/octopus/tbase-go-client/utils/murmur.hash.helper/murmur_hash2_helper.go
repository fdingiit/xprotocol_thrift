package murmur_hash_helper

import (
	"unsafe"
)

/*
 * -----------------------------------------------------------------------------
 * MurmurHash2, by Austin Appleby
 * Note - This code makes a few assumptions about how your machine behaves -
 * 1. We can read a 4-byte value from any address without crashing
 * 2. sizeof(int) == 4
 * And it has a few limitations -
 * 1. It will not work incrementally.
 * 2. It will not produce the same results on little-endian and big-endian
 *    machines.
 */

func MurmurHash2(key []byte, seed int32) int64 {

	m := ConvertToInt32(0x5bd1e995)
	r := uint32(24)
	dataLen := uint32(len(key))
	len := dataLen
	/* Initialize the hash to a 'random' value */
	h := seed ^ int32(dataLen)

	/* Mix 4 bytes at a time into the hash */
	data := key

	i := 0
	for {
		if dataLen < 4 {
			break
		}
		k := *(*int32)(unsafe.Pointer(&data[i*4]))

		k *= m
		k ^= int32(uint32(k) >> r)
		k *= m

		h *= m
		h ^= k

		i += 1
		dataLen -= 4
	}

	/* Handle the last few bytes of the input array */
	switch len % 4 {
	case 3:
		h ^= int32(data[(len&^3)+2]&0xff) << 16
		h ^= int32(data[(len&^3)+1]&0xff) << 8
		h ^= int32(data[(len &^ 3)] & 0xff)
		h *= m
	case 2:
		h ^= int32(data[(len&^3)+1]&0xff) << 8
		h ^= int32(data[(len &^ 3)] & 0xff)
		h *= m
	case 1:
		h ^= int32(data[(len &^ 3)] & 0xff)
		h *= m
	default:
	}

	h ^= int32(uint32(h) >> 13)
	h *= m
	h ^= int32(uint32(h) >> 15)

	return int64(h) & 0x00000000ffffffff
}

func ConvertToInt32(data int) int32 {
	i64 := int64(data)
	return int32(i64)
}
