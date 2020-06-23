package util

import "encoding/binary"

// convert uint64 to 8-byte big endian representation
func Utob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
