package deserializer

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	curve "github.com/consensys/gnark-crypto/ecc/bn254"
)

type Phase1 struct {
	tauG1      []curve.G1Affine
	alphaTauG1 []curve.G1Affine
	betaTauG1  []curve.G1Affine
	tauG2      []curve.G2Affine
	betaG2     curve.G2Affine
}

func ConvertPtauToPhase1(ptau Ptau) (phase1 Phase1, err error) {
	tauG1 := make([]curve.G1Affine, len(ptau.PTauPubKey.TauG1))
	for i, g1 := range ptau.PTauPubKey.TauG1 {
		g1Affine := curve.G1Affine{}
		x := bytesToElement(g1[0].Bytes())
		g1Affine.X = x
		y := bytesToElement(g1[1].Bytes())
		g1Affine.Y = y
		// fmt.Printf("X: %v \n", g1Affine.X.String())
		// fmt.Printf("Y: %v \n", g1Affine.Y.String())
		// fmt.Printf("g1Affine: %v \n", g1Affine)
		if !g1Affine.IsOnCurve() {
			fmt.Errorf("tauG1: \n index: %v g1Affine.X: %v \n g1Affine.Y: %v \n", i, g1Affine.X.String(), g1Affine.Y.String())
			panic("g1Affine is not on curve")
		}
		tauG1[i] = g1Affine
	}

	alphaTauG1 := make([]curve.G1Affine, len(ptau.PTauPubKey.AlphaTauG1))
	for i, g1 := range ptau.PTauPubKey.AlphaTauG1 {
		g1Affine := curve.G1Affine{}
		x := bytesToElement(g1[0].Bytes())
		g1Affine.X = x
		y := bytesToElement(g1[1].Bytes())
		g1Affine.Y = y
		if !g1Affine.IsOnCurve() {
			fmt.Errorf("alphaTauG1: \n index: %v g1Affine.X: %v \n g1Affine.Y: %v \n", i, g1Affine.X.String(), g1Affine.Y.String())
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
			fmt.Errorf("betaTauG1: \n index: %v, g1Affine.X: %v \n g1Affine.Y: %v \n", i, g1Affine.X.String(), g1Affine.Y.String())
			panic("g1Affine is not on curve")
		}
		betaTauG1[i] = g1Affine
	}
	// fmt.Printf("betaTauG1: %v \n", betaTauG1)

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

		// fmt.Printf("X: %v \n", g2Affine.X.String())
		// fmt.Printf("Y: %v \n", g2Affine.Y.String())
		// fmt.Printf("g2Affine %v: %v \n", i, g2Affine)
		if !g2Affine.IsOnCurve() {
			fmt.Printf("tauG2: \n index: %v, g2Affine.X.A0: %v \n g2Affine.X.A1: %v \n g2Affine.Y.A0: %v \n g2Affine.Y.A1 %v \n", i, g2Affine.X.A0.String(), g2Affine.X.A1.String(), g2Affine.Y.A0.String(), g2Affine.Y.A1.String())
			panic("g2Affine is not on curve")
		}
		tauG2[i] = g2Affine
	}

	// fmt.Printf("tauG2: %v \n", tauG2)

	betaG2 := curve.G2Affine{}
	{
		g2 := ptau.PTauPubKey.BetaG2

		x0 := bytesToElement(g2[0].Bytes())
		x1 := bytesToElement(g2[1].Bytes())
		betaG2.X.A0 = x0
		betaG2.X.A1 = x1
		y0 := bytesToElement(g2[2].Bytes())
		y1 := bytesToElement(g2[3].Bytes())
		betaG2.Y.A0 = y0
		betaG2.Y.A1 = y1

		if !betaG2.IsOnCurve() {
			fmt.Printf("g2Affine.X.A0: %v \n g2Affine.X.A1: %v \n g2Affine.Y.A0: %v \n g2Affine.Y.A1 %v \n", betaG2.X.A0.String(), betaG2.X.String(), betaG2.Y.A0.String(), betaG2.Y.A1.String())
			panic("g2Affine is not on curve")
		}
	}

	//fmt.Printf("BetaG2: %v \n", BetaG2)

	return Phase1{tauG1: tauG1, tauG2: tauG2, alphaTauG1: alphaTauG1, betaTauG1: betaTauG1, betaG2: betaG2}, nil
}

func WritePhase1FromPtauFile(ptauFile *PtauFile, outputPath string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var header Header

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	N := ptauFile.DomainSize()

	fmt.Printf("Power %d supports up to %d constraints\n", ptauFile.Header.Power, N)

	header.Power = byte(ptauFile.Header.Power)

	// can be extracted from ptau.Contributions (7) but hardcoding for now
	// ptau link: https://github.com/iden3/snarkjs/tree/master#7-prepare-phase-2
	header.Contributions = 54

	// Write the header
	err = header.writeTo(outputFile)
	if err != nil {
		return err
	}

	// BN254 encoder using compressed representation of points to save storage space
	enc := bn254.NewEncoder(writer)
	fmt.Println("1. Writing TauG1")
	tauG1 := make(chan curve.G1Affine, 10000)
	go ptauFile.ReadTauG1(tauG1)
	for point := range tauG1 {
		if err := enc.Encode(&point); err != nil {
			return err
		}
	}

	// Write α[τ⁰]₁, α[τ¹]₁, α[τ²]₁, …, α[τᴺ⁻¹]₁
	fmt.Println("2. Writing AlphaTauG1")
	alphaTauG1 := make(chan curve.G1Affine, 10000)
	go ptauFile.ReadAlphaTauG1(alphaTauG1)
	for point := range alphaTauG1 {
		if err := enc.Encode(&point); err != nil {
			return err
		}
	}

	// Write β[τ⁰]₁, β[τ¹]₁, β[τ²]₁, …, β[τᴺ⁻¹]₁
	fmt.Println("3. Writing BetaTauG1")
	betaTauG1 := make(chan curve.G1Affine, 10000)
	go ptauFile.ReadBetaTauG1(betaTauG1)
	for point := range betaTauG1 {
		if err := enc.Encode(&point); err != nil {
			return err
		}
	}

	// Write {[τ⁰]₂, [τ¹]₂, [τ²]₂, …, [τᴺ⁻¹]₂}
	fmt.Println("4. Writing TauG2")
	tauG2 := make(chan curve.G2Affine, 10000)
	go ptauFile.ReadTauG2(tauG2)
	for point := range tauG2 {
		if err := enc.Encode(&point); err != nil {
			return err
		}
	}

	// Write [β]₂
	fmt.Println("5. Writing BetaG2")
	betaG2, err := ptauFile.ReadBetaG2()
	if err != nil {
		return err
	}
	enc.Encode(&betaG2)
	return nil
}

func WritePhase1(phase1 Phase1, power byte, outputPath string) error {
	// output outputFile
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var header Header

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	N := int(math.Pow(2, float64(power)))

	fmt.Printf("Power %d supports up to %d constraints\n", power, N)

	header.Power = power

	// can be extracted from ptau.Contributions (7) but hardcoding for now
	// ptau link: https://github.com/iden3/snarkjs/tree/master#7-prepare-phase-2
	header.Contributions = 54

	// Write the header
	header.writeTo(outputFile)

	// BN254 encoder using compressed representation of points to save storage space
	enc := bn254.NewEncoder(writer)

	// Taken from https://github.com/worldcoin/semaphore-mtb-setup/blob/main/phase1/phase1.go
	// In the initialization, τ = α = β = 1, so we are writing the generators directly
	// Write [τ⁰]₁, [τ¹]₁, [τ²]₁, …, [τ²ᴺ⁻²]₁
	fmt.Println("1. Writing TauG1")
	for i := 0; i < 2*N-1; i++ {
		if err := enc.Encode(&phase1.tauG1[i]); err != nil {
			return err
		}
	}

	// Write α[τ⁰]₁, α[τ¹]₁, α[τ²]₁, …, α[τᴺ⁻¹]₁
	fmt.Println("2. Writing AlphaTauG1")
	for i := 0; i < N; i++ {
		if err := enc.Encode(&phase1.alphaTauG1[i]); err != nil {
			return err
		}
	}

	// Write β[τ⁰]₁, β[τ¹]₁, β[τ²]₁, …, β[τᴺ⁻¹]₁
	fmt.Println("3. Writing BetaTauG1")
	for i := 0; i < N; i++ {
		if err := enc.Encode(&phase1.betaTauG1[i]); err != nil {
			return err
		}
	}

	// Write {[τ⁰]₂, [τ¹]₂, [τ²]₂, …, [τᴺ⁻¹]₂}
	fmt.Println("4. Writing TauG2")
	for i := 0; i < N; i++ {
		if err := enc.Encode(&phase1.tauG2[i]); err != nil {
			return err
		}
	}

	// Write [β]₂
	fmt.Println("5. Writing BetaG2")
	enc.Encode(&phase1.betaG2)

	return nil
}
