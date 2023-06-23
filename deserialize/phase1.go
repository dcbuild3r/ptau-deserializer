package deserializer

import (
	"encoding/binary"
	"fmt"

	curve "github.com/consensys/gnark-crypto/ecc/bn254"
	fp "github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark/backend/groth16/bn254/mpcsetup"
)

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

func convertPtauToSrs(ptau Ptau) (phase1 mpcsetup.Phase1, err error) {
	tauG1 := make([]curve.G1Affine, len(ptau.PTauPubKey.TauG1))
	for i, g1 := range ptau.PTauPubKey.TauG1 {
		g1Affine := curve.G1Affine{}
		x := bytesToElement(g1[0].Bytes())
		g1Affine.X = x
		y := bytesToElement(g1[1].Bytes())
		g1Affine.Y = y
		fmt.Printf("X: %v \n", g1Affine.X.String())
		fmt.Printf("Y: %v \n", g1Affine.Y.String())
		fmt.Printf("g1Affine: %v \n", g1Affine)
		if !g1Affine.IsOnCurve() {
			panic("g1Affine is not on curve")
		}
		tauG1[i] = g1Affine
	}

	// fmt.Printf("tauG1: %v \n", tauG1)

	tauG2 := make([]curve.G2Affine, len(ptau.PTauPubKey.TauG2))
	for i, g2 := range ptau.PTauPubKey.TauG2 {
		g2Affine := curve.G2Affine{}
		x0 := bytesToElement(g2[0].Bytes())
		x1 := bytesToElement(g2[1].Bytes())
		g2Affine.X.A0 = x0
		g2Affine.X.A1 = x1
		y0 := bytesToElement(g2[2].Bytes())
		y1 := bytesToElement(g2[3].Bytes())
		g2Affine.Y.A0 = y0
		g2Affine.Y.A1 = y1

		fmt.Printf("X: %v \n", g2Affine.X.String())
		fmt.Printf("Y: %v \n", g2Affine.Y.String())
		fmt.Printf("g2Affine %v: %v \n", i, g2Affine)
		if !g2Affine.IsOnCurve() {
			panic("g2Affine is not on curve")
		}
		tauG2[i] = g2Affine
	}

	// fmt.Printf("tauG1: %v \n", tauG1)

	alphaTauG1 := make([]curve.G1Affine, len(ptau.PTauPubKey.AlphaTauG1))
	for i, g1 := range ptau.PTauPubKey.AlphaTauG1 {
		g1Affine := curve.G1Affine{}
		x := bytesToElement(g1[0].Bytes())
		g1Affine.X = x
		y := bytesToElement(g1[1].Bytes())
		g1Affine.Y = y
		if !g1Affine.IsOnCurve() {
			panic("g1Affine is not on curve")
		}
		alphaTauG1[i] = g1Affine
	}
	// fmt.Printf("alphaTauG1: %v \n", alphaTauG1)

	betaTauG1 := make([]curve.G1Affine, len(ptau.PTauPubKey.BetaTauG1))

	for i, g1 := range ptau.PTauPubKey.BetaTauG1 {
		g1Affine := curve.G1Affine{}
		x := bytesToElement(g1[0].Bytes())
		g1Affine.X = x
		y := bytesToElement(g1[1].Bytes())
		g1Affine.Y = y
		if !g1Affine.IsOnCurve() {
			panic("g1Affine is not on curve")
		}
		betaTauG1[i] = g1Affine
	}
	// fmt.Printf("betaTauG1: %v \n", betaTauG1)

	TauG2 := make([]curve.G2Affine, len(ptau.PTauPubKey.TauG2))

	for i, g2 := range ptau.PTauPubKey.TauG2 {
		g2Affine := curve.G2Affine{}
		x0 := bytesToElement(g2[0].Bytes())
		x1 := bytesToElement(g2[1].Bytes())
		g2Affine.X.A0 = x0
		g2Affine.X.A1 = x1
		y0 := bytesToElement(g2[2].Bytes())
		y1 := bytesToElement(g2[3].Bytes())
		g2Affine.Y.A0 = y0
		g2Affine.Y.A1 = y1

		if !g2Affine.IsOnCurve() {
			panic("g2Affine is not on curve")
		}
		TauG2[i] = g2Affine
	}

	//fmt.Printf("TauG2: %v \n", TauG2)

	BetaG2 := curve.G2Affine{}

	{
		g2 := ptau.PTauPubKey.BetaG2

		x0 := bytesToElement(g2[0].Bytes())
		x1 := bytesToElement(g2[1].Bytes())
		BetaG2.X.A0 = x0
		BetaG2.X.A1 = x1
		y0 := bytesToElement(g2[2].Bytes())
		y1 := bytesToElement(g2[3].Bytes())
		BetaG2.Y.A0 = y0
		BetaG2.Y.A1 = y1

		if !BetaG2.IsOnCurve() {
			panic("g2Affine is not on curve")
		}
	}

	//fmt.Printf("BetaG2: %v \n", BetaG2)

	phase1.Parameters.G1.Tau = tauG1
	phase1.Parameters.G1.AlphaTau = alphaTauG1
	phase1.Parameters.G1.BetaTau = betaTauG1
	phase1.Parameters.G2.Tau = TauG2
	phase1.Parameters.G2.Beta = BetaG2

	return phase1, nil
}
