package IPTools

import (
	"strconv"
	"strings"
)

func IPv4ToNumberForm(ipv4 string) uint32 {
	splitResult := strings.Split(ipv4, ".")
	var numberForm uint32 = 0
	for _, s := range splitResult {
		num, _ := strconv.ParseUint(s, 10, 32)
		numberForm <<= 8
		numberForm |= uint32(num)
	}
	return numberForm
}

func makeMask(maskBitLength uint8) uint32 {
	result := uint32(0)
	for i := uint8(0); i <= maskBitLength; i++ {
		result |= 1 << (32 - i)
	}
	return result
}

func MaskIP(ip string, maskBitLength uint8) (masked uint32) {
	return IPv4ToNumberForm(ip) & makeMask(maskBitLength)
}

func jhash(arr []uint32) string {
	var h uint32 = 8388617
	l := len(arr)
	for i := 0; i < l; i++ {
		h = ((h<<1 | h>>30) & 0x7fffffff) ^ arr[i]
	}
	return strconv.FormatUint(
		uint64(h), 36,
	)
}

func JHashWithMask(ip string, maskBitLength uint8) string {
	return jhash([]uint32{MaskIP(ip, maskBitLength), uint32(maskBitLength)})
}
