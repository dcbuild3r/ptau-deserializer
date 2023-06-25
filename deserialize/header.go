package deserializer

import (
	"encoding/binary"
	"io"
)

type Header struct {
	Power         byte
	Contributions uint16
}

func (p *Header) ReadFrom(reader io.Reader) error {
	buffPower := make([]byte, 1)
	// Read NConstraints
	if _, err := reader.Read(buffPower); err != nil {
		return err
	}
	p.Power = buffPower[0]

	// Read NContribution
	buffContributions := make([]byte, 2)
	if _, err := reader.Read(buffContributions); err != nil {
		return err
	}
	p.Contributions = binary.BigEndian.Uint16(buffContributions)
	return nil
}

func (p *Header) writeTo(writer io.Writer) error {
	// Write Power
	if _, err := writer.Write([]byte{p.Power}); err != nil {
		return err
	}

	// Write Contribution
	buffContributions := make([]byte, 2)
	binary.BigEndian.PutUint16(buffContributions, p.Contributions)
	if _, err := writer.Write(buffContributions); err != nil {
		return err
	}

	return nil
}
