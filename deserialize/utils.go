package deserializer

import (
	"encoding/binary"
	"io"
	"math/big"

	fp "github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

func readULE32(reader io.Reader) (uint32, error) {
	var buffer = make([]byte, 4)

	_, err := reader.Read(buffer)

	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(buffer), nil
}

func readULE64(reader io.Reader) (uint64, error) {
	var buffer = make([]byte, 8)

	_, err := reader.Read(buffer)

	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(buffer), nil
}

func readBigInt(reader io.Reader, n8 uint32) (big.Int, error) {
	var buffer = make([]byte, n8)

	_, err := reader.Read(buffer)
	reverseSlice(buffer)

	if err != nil {
		return *big.NewInt(0), err
	}

	bigInt := big.NewInt(0).SetBytes(buffer)

	return *bigInt, nil
}

func reverseSlice(slice []byte) []byte {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func bytesToElement(b []byte) fp.Element {
	var z fp.Element
	reverseSlice(b)
	if len(b) < 32 {
		b = append(b, make([]byte, 32-len(b))...)
	}

	z[0] = binary.LittleEndian.Uint64(b[0:8])
	z[1] = binary.LittleEndian.Uint64(b[8:16])
	z[2] = binary.LittleEndian.Uint64(b[16:24])
	z[3] = binary.LittleEndian.Uint64(b[24:32])

	return z
}
