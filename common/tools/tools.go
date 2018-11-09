package tools

func PackUInt32(a uint32, b uint32) uint64 {
	return uint64(a)<<32 | uint64(b)
}
