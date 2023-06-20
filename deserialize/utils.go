package deserializer

import (
	"bufio"
	"encoding/binary"
	"math/big"
)

func readULE32(reader *bufio.Reader) (uint32, error) {
	var buffer = make([]byte, 4)

	_, err := reader.Read(buffer)

	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(buffer), nil
}

func readBigInt(reader *bufio.Reader, n8 uint32) (big.Int, error) {
	var buffer = make([]byte, n8)

	_, err := reader.Read(buffer)

	if err != nil {
		return *big.NewInt(0), err
	}

	bigInt := big.NewInt(int64(binary.LittleEndian.Uint64(buffer)))

	return *bigInt, nil
}
