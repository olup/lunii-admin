/**
This is a *custom* implementation of the xxtea algorithm.
*/
package lunii

import (
	"bytes"
	"encoding/binary"
)

const delta = 0x9e3779b9

func mx(y, z, p, e, sum uint32, key []uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (key[(p&3)^e] ^ z))
}

func Btea(v []uint32, n int, key []uint32) []uint32 {
	var i, y, z, p, e, sum uint32
	// Coding Part
	if n > 1 {
		un := uint32(n)
		rounds := 1 + 52/un

		z = v[n-1]

		for {
			sum += delta
			e = (sum >> 2) & 3
			for p = 0; p < un-1; p++ {
				y = v[p+1]
				v[p] += mx(y, z, p, e, sum, key)
				z = v[p]
			}

			y = v[0]
			v[n-1] += mx(y, z, p, e, sum, key)
			z = v[n-1]

			i++
			if i > rounds-1 {
				break
			}
		}

	} else if n < -1 { // Decoding Part
		un := uint32(-n)
		rounds := 1 + 52/un

		sum = rounds * delta
		y = v[0]

		for {
			e = (sum >> 2) & 3
			for p = un - 1; p > 0; p-- {
				z = v[p-1]
				v[p] -= mx(y, z, p, e, sum, key)
				y = v[p]
			}

			z = v[un-1]
			v[0] -= mx(y, z, p, e, sum, key)
			y = v[0]
			sum -= delta

			i++
			if i > rounds-1 {
				break
			}
		}
	}
	return v
}

func bytesToInt32(in []byte, endianness binary.ByteOrder) []uint32 {
	var out = make([]uint32, len(in)/4)
	reader := bytes.NewReader(in)
	for i := 0; i < len(in)/4; i++ {
		var result uint32
		binary.Read(reader, endianness, &result)
		out[i] = result
	}
	return out
}

func int32sToBytes(in []uint32, endianness binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)

	for i := 0; i < len(in); i++ {
		binary.Write(buf, endianness, &in[i])
	}
	out := buf.Bytes()
	return out
}
