package deserializer

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
)

const GROTH_16_PROTOCOL_ID = 1

type NotGroth16 struct {
	Err error
}

func (r *NotGroth16) Error() string {
	return fmt.Sprintf("Groth16 is the only supported protocol at this time (PLONK and FFLONK are not): %v", r.Err)
}

// Incomplete (only extracts necessary fields for conversion to .ph1 format)
type Zkey struct {
	Header         Header
	protocolHeader HeaderGroth
}

type Header struct {
	ProtocolID     uint32
	protocolHeader HeaderGroth
}

type HeaderGroth struct {
	n8q        uint32
	q          big.Int
	n8r        uint32
	r          big.Int
	nVars      uint32
	nPublic    uint32
	domainSize uint32
	power      float64
}

func ReadZkey(zkeyPath string) (Zkey, error) {
	file, err := os.Open(zkeyPath)

	if err != nil {
		return Zkey{}, err
	}

	defer file.Close()

	// Create a new buffered reader
	reader := bufio.NewReader(file)

	header, err := readHeader(reader)

	if err != nil {
		return Zkey{}, err
	}

	zkey := Zkey{Header: header, protocolHeader: header.protocolHeader}

	return zkey, nil
}

func readHeader(reader *bufio.Reader) (Header, error) {
	var header = Header{}

	protocolID, err := readULE32(reader)

	if err != nil {
		return header, err
	}

	// if groth16
	if protocolID == GROTH_16_PROTOCOL_ID {
		headerGroth, err := readHeaderGroth16(reader)

		if err != nil {
			return header, err
		}

		header = Header{ProtocolID: protocolID, protocolHeader: headerGroth}

	} else {
		return header, &NotGroth16{Err: errors.New("ProtocolID is not Groth16")}
	}

	return header, nil
}

func readHeaderGroth16(reader *bufio.Reader) (HeaderGroth, error) {
	var header = HeaderGroth{}

	n8q, err := readULE32(reader)

	if err != nil {
		return header, err
	}

	q, err := readBigInt(reader, n8q)

	if err != nil {
		return header, err
	}

	n8r, err := readULE32(reader)

	if err != nil {
		return header, err
	}

	r, err := readBigInt(reader, n8r)

	if err != nil {
		return header, err
	}

	nVars, err := readULE32(reader)

	if err != nil {
		return header, err
	}

	nPublic, err := readULE32(reader)

	if err != nil {
		return header, err
	}

	domainSize, err := readULE32(reader)

	if err != nil {
		return header, err
	}

	power := math.Log2(float64(domainSize))

	header = HeaderGroth{n8q: n8q, q: q, n8r: n8r, r: r, nVars: nVars, nPublic: nPublic, domainSize: domainSize, power: power}

	return header, nil
}
