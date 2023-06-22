package deserializer

import (
	"fmt"
	"io"
	"math/big"
	"os"
)

///////////////////////////////////////////////////////////////////
///                             PTAU                            ///
///////////////////////////////////////////////////////////////////

// Format
// Taken from the iden3/snarkjs repo powersoftau_new.js file
// https://github.com/iden3/snarkjs/blob/master/src/powersoftau_new.js
/*
Header(1)
    n8
    prime
    power
tauG1(2)
    {(2 ** power)*2-1} [
        G1, tau*G1, tau^2 * G1, ....
    ]
tauG2(3)
    {2 ** power}[
        G2, tau*G2, tau^2 * G2, ...
    ]
alphaTauG1(4)
    {2 ** power}[
        alpha*G1, alpha*tau*G1, alpha*tau^2*G1,....
    ]
betaTauG1(5)
    {2 ** power} []
        beta*G1, beta*tau*G1, beta*tau^2*G1, ....
    ]
betaG2(6)
    {1}[
        beta*G2
    ]
contributions(7) - Ignore contributions, users can verify using snarkjs
    NContributions
    {NContributions}[
        tau*G1
        tau*G2
        alpha*G1
        beta*G1
        beta*G2
        pubKey
            tau_g1s
            tau_g1sx
            tau_g2spx
            alpha_g1s
            alpha_g1sx
            alpha_g1spx
            beta_g1s
            beta_g1sx
            beta_g1spx
        partialHash (216 bytes) See https://github.com/mafintosh/blake2b-wasm/blob/23bee06945806309977af802bc374727542617c7/blake2b.wat#L9
        hashNewChallenge
    ]
*/

// in bytes
const BN254_FIELD_ELEMENT_SIZE = 32

// G1 and G2 are both arrays of two big.Ints (field elements)
type G1 [2]big.Int
type G2 [4]big.Int

type PtauHeader struct {
	n8    uint32
	prime big.Int
	power uint32
}

type Ptau struct {
	Header     PtauHeader
	PTauPubKey PtauPubKey
}

type PtauPubKey struct {
	TauG1      []G1
	TauG2      []G2
	AlphaTauG1 []G1
	BetaTauG1  []G1
	BetaG2     G2
}

func ReadPtau(zkeyPath string) (Ptau, error) {
	reader, err := os.Open(zkeyPath)

	if err != nil {
		return Ptau{}, err
	}

	defer reader.Close()

	var ptauStr = make([]byte, 4)
	_, err = reader.Read(ptauStr)

	fmt.Printf("zkeyStr: %s \n", string(ptauStr))

	// version
	_, err = readULE32(reader)

	// number of sections
	_, err = readULE32(reader)

	numSections := uint32(7)
	fmt.Printf("num sections: %v \n", numSections)

	// in practice, all sections have only one segment, but who knows...
	// 1-based indexing, so we need to allocate one more than the number of sections
	sections := make([][]SectionSegment, numSections+1)
	for i := uint32(0); i < numSections; i++ {
		ht, _ := readULE32(reader)
		hl, _ := readULE64(reader)
		fmt.Printf("ht: %v \n", ht)
		fmt.Printf("hl: %v \n", hl)
		if sections[ht] == nil {
			sections[ht] = make([]SectionSegment, 0)
		}
		pos, _ := reader.Seek(0, io.SeekCurrent)
		sections[ht] = append(sections[ht], SectionSegment{pos: uint32(pos), size: hl})
		reader.Seek(int64(hl), io.SeekCurrent)
	}

	fmt.Printf("sections: %v \n", sections)

	// section size
	_, err = readBigInt(reader, 8)

	// Header (1)
	seekToUniqueSection(reader, sections, 1)

	// Read header
	header, err := readPtauHeader(reader)

	if err != nil {
		return Ptau{}, err
	}

	// TauG1 (2)
	seekToUniqueSection(reader, sections, 2)

	var PtauPubKey PtauPubKey

	twoToPower := uint32(1 << header.power)
	PtauPubKey.TauG1, err = readG1Array(reader, twoToPower*2-1)

	if err != nil {
		return Ptau{}, err
	}

	// TauG2 (3)
	seekToUniqueSection(reader, sections, 3)

	PtauPubKey.TauG2, err = readG2Array(reader, twoToPower)

	if err != nil {
		return Ptau{}, err
	}

	// AlphaTauG1 (4)
	seekToUniqueSection(reader, sections, 4)

	PtauPubKey.AlphaTauG1, err = readG1Array(reader, twoToPower)

	if err != nil {
		return Ptau{}, err
	}

	// BetaTauG1 (5)
	seekToUniqueSection(reader, sections, 5)

	PtauPubKey.BetaTauG1, err = readG1Array(reader, twoToPower)

	if err != nil {
		return Ptau{}, err
	}

	// BetaG2 (6)
	seekToUniqueSection(reader, sections, 6)

	PtauPubKey.BetaG2, err = readG2(reader)

	if err != nil {
		return Ptau{}, err
	}

	return Ptau{Header: header, PTauPubKey: PtauPubKey}, nil
}

// looks good ser
// nice!

// thanks! btw now I have to just figure out how to serialize it into .ph1 and then just write

// don't think so - can call the function directly

func readPtauHeader(reader io.ReadSeeker) (PtauHeader, error) {
	var header PtauHeader

	n8, err := readULE32(reader)

	if err != nil {
		return PtauHeader{}, err
	}

	header.n8 = n8

	prime, err := readBigInt(reader, n8)

	if err != nil {
		return PtauHeader{}, err
	}

	header.prime = prime

	power, err := readULE32(reader)

	if err != nil {
		return PtauHeader{}, err
	}

	header.power = power

	return header, nil
}

func readG1Array(reader io.ReadSeeker, numPoints uint32) ([]G1, error) {
	g1s := make([]G1, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		g1, err := readG1(reader)

		if err != nil {
			return []G1{}, err
		}

		g1s[i] = g1
	}
	return g1s, nil
}

func readG2Array(reader io.ReadSeeker, numPoints uint32) ([]G2, error) {
	g2s := make([]G2, numPoints)

	for i := uint32(0); i < numPoints; i++ {
		g2, err := readG2(reader)

		if err != nil {
			return []G2{}, err
		}

		g2s[i] = g2
	}

	return g2s, nil
}

func readTauG2(reader io.ReadSeeker) ([]G2, error) {
	tauG2_s, err := readG2(reader)

	if err != nil {
		return []G2{}, err
	}

	tauG2_sx, err := readG2(reader)

	if err != nil {
		return []G2{}, err
	}

	return []G2{tauG2_s, tauG2_sx}, nil
}

func readG1(reader io.ReadSeeker) (G1, error) {
	var g1 G1

	x, err := readBigInt(reader, BN254_FIELD_ELEMENT_SIZE)

	if err != nil {
		return G1{}, err
	}

	g1[0] = x

	y, err := readBigInt(reader, BN254_FIELD_ELEMENT_SIZE)

	if err != nil {
		return G1{}, err
	}

	g1[1] = y

	return g1, nil
}

func readG2(reader io.ReadSeeker) (G2, error) {
	var g2 G2

	x0, err := readBigInt(reader, BN254_FIELD_ELEMENT_SIZE)

	if err != nil {
		return G2{}, err
	}

	g2[0] = x0

	x1, err := readBigInt(reader, BN254_FIELD_ELEMENT_SIZE)

	if err != nil {
		return G2{}, err
	}

	g2[1] = x1

	y0, err := readBigInt(reader, BN254_FIELD_ELEMENT_SIZE)

	if err != nil {
		return G2{}, err
	}

	g2[2] = y0

	y1, err := readBigInt(reader, BN254_FIELD_ELEMENT_SIZE)

	if err != nil {
		return G2{}, err
	}

	g2[3] = y1

	return g2, nil
}
