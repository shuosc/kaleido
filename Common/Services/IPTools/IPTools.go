package IPTools

import (
	"net"
	"strconv"
	"strings"
)

func IPv4ToNumberForm(ipv4 string) (uint32, error) {
	splitResult := strings.Split(ipv4, ".")
	var numberForm uint32 = 0
	for _, s := range splitResult {
		num, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return 0, &net.AddrError{Err: "Not a valid ipv4 Address", Addr: ipv4}
		}
		numberForm <<= 8
		numberForm |= uint32(num)
	}
	return numberForm, nil
}

func makeMask(maskBitLength uint8) uint32 {
	result := uint32(0)
	for i := uint8(0); i <= maskBitLength; i++ {
		result |= 1 << (32 - i)
	}
	return result
}

var maskTable [33]uint32

func MaskIP(ipNumberForm uint32, maskBitLength uint8) (masked uint32) {
	return ipNumberForm & maskTable[maskBitLength]
}

func init() {
	for i := uint8(0); i <= 32; i++ {
		maskTable[i] = makeMask(i)
	}
}
