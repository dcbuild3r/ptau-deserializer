package deserializer

import (
	"fmt"

	curve "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark/backend/groth16/bn254/mpcsetup"
)

func convertPtauToSrs(ptau Ptau) (mpcsetup.Phase1, error) {
	tauG1 := make([]curve.G1Affine, len(ptau.PTauPubKey.TauG1))
	for i, g1 := range ptau.PTauPubKey.TauG1 {
		g1Affine := curve.G1Affine{}
		g1Affine.X.SetBigInt(&g1[0])
		g1Affine.Y.SetBigInt(&g1[1])
		tauG1[i] = g1Affine
	}
	fmt.Printf("tauG1: %v \n", tauG1)
	return mpcsetup.Phase1{}, nil
}
