package base36

import (
	"math"
)

// base36 implements by github.com/martinlindhe/base36
// tinyurl only returns lowercase but requests are case insensitive

var (
	base36 = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z'}

	index = map[byte]int{
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4,
		'5': 5, '6': 6, '7': 7, '8': 8, '9': 9,

		'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14,
		'f': 15, 'g': 16, 'h': 17, 'i': 18, 'j': 19,
		'k': 20, 'l': 21, 'm': 22, 'n': 23, 'o': 24,
		'p': 25, 'q': 26, 'r': 27, 's': 28, 't': 29,
		'u': 30, 'v': 31, 'w': 32, 'x': 33, 'y': 34,
		'z': 35,

		'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14,
		'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19,
		'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24,
		'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29,
		'U': 30, 'V': 31, 'W': 32, 'X': 33, 'Y': 34,
		'Z': 35,
	}
)

// Encode encodes a uint64 to base36 string.
// the max length is 8 bytes (8 chars)
func Encode(value uint64) string {
	var res [8]byte
	var i int
	for i = len(res) - 1; ; i-- {
		res[i] = base36[value%36]
		value /= 36
		if value == 0 {
			break
		}
	}

	return string(res[i:])
}

// Decode decodes a base36-encoded string to uint64.
func Decode(s string) uint64 {
	res := uint64(0)
	l := len(s) - 1
	for idx := range s {
		c := s[l-idx]
		res += uint64(index[c]) * uint64(math.Pow(36, float64(idx)))
	}
	return res
}
