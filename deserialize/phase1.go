package deserializer

import (
	"encoding/binary"
	"fmt"
	"io"

	curve "github.com/consensys/gnark-crypto/ecc/bn254"
	fp "github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

// Taken from https://github.com/ConsenSys/gnark/blob/b01ec42ceb18e1d316e3b7386e195b929212325e/backend/groth16/bn254/mpcsetup/phase1.go
type Phase1 struct {
	Parameters struct {
		G1 struct {
			Tau      []curve.G1Affine // {[τ⁰]₁, [τ¹]₁, [τ²]₁, …, [τ²ⁿ⁻²]₁}
			AlphaTau []curve.G1Affine // {α[τ⁰]₁, α[τ¹]₁, α[τ²]₁, …, α[τⁿ⁻¹]₁}
			BetaTau  []curve.G1Affine // {β[τ⁰]₁, β[τ¹]₁, β[τ²]₁, …, β[τⁿ⁻¹]₁}
		}
		G2 struct {
			Tau  []curve.G2Affine // {[τ⁰]₂, [τ¹]₂, [τ²]₂, …, [τⁿ⁻¹]₂}
			Beta curve.G2Affine   // [β]₂
		}
	}
	PublicKeys struct {
		Tau, Alpha, Beta PublicKey
	}
	Hash []byte // sha256 hash
}

// taken from https://github.com/ConsenSys/gnark/blob/b01ec42ceb18e1d316e3b7386e195b929212325e/backend/groth16/bls12-381/mpcsetup/marshal.go
// WriteTo implements io.WriterTo
func (phase1 *Phase1) WriteTo(writer io.Writer) (int64, error) {
	n, err := phase1.writeTo(writer)
	if err != nil {
		return n, err
	}
	nBytes, err := writer.Write(phase1.Hash)
	return int64(nBytes) + n, err
}

func (phase1 *Phase1) writeTo(writer io.Writer) (int64, error) {
	toEncode := []interface{}{
		&phase1.PublicKeys.Tau.SG,
		&phase1.PublicKeys.Tau.SXG,
		&phase1.PublicKeys.Tau.XR,
		&phase1.PublicKeys.Alpha.SG,
		&phase1.PublicKeys.Alpha.SXG,
		&phase1.PublicKeys.Alpha.XR,
		&phase1.PublicKeys.Beta.SG,
		&phase1.PublicKeys.Beta.SXG,
		&phase1.PublicKeys.Beta.XR,
		phase1.Parameters.G1.Tau,
		phase1.Parameters.G1.AlphaTau,
		phase1.Parameters.G1.BetaTau,
		phase1.Parameters.G2.Tau,
		&phase1.Parameters.G2.Beta,
	}

	enc := curve.NewEncoder(writer)
	for _, v := range toEncode {
		if err := enc.Encode(v); err != nil {
			return enc.BytesWritten(), err
		}
	}
	return enc.BytesWritten(), nil
}

// taken from https://github.com/ConsenSys/gnark/blob/b01ec42ceb18e1d316e3b7386e195b929212325e/backend/groth16/bls12-381/mpcsetup/utils.go
type PublicKey struct {
	SG  curve.G1Affine
	SXG curve.G1Affine
	XR  curve.G2Affine
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

func convertPtauToSrs(ptau Ptau) (phase1 Phase1, err error) {
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
