package deserializer

import (
	"encoding/binary"
	"io"
	"math/big"
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
