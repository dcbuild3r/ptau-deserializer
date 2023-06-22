package deserializer

import "math/big"

///////////////////////////////////////////////////////////////////
///                             PPOT                            ///
///////////////////////////////////////////////////////////////////

// Taken from the iden3/snarkjs repo powersoftau_new.js file
// https://github.com/iden3/snarkjs/blob/fb144555d8ce4779ad79e707f269771c672a8fb7/src/powersoftau_challenge_contribute.js#L45
// Format of the output
//      Hash of the last contribution  64 Bytes
//      2^N*2-1 TauG1 Points (compressed)
//      2^N TauG2 Points (compressed)
//      2^N AlphaTauG1 Points (compressed)
//      2^N BetaTauG1 Points (compressed)
//      Public Key
//          BetaG2 (compressed)
//          G1*s (compressed)
//          G1*s*tau (compressed)
//          G1*t (compressed)
//          G1*t*alpha (compressed)
//          G1*u (compressed)
//          G1*u*beta (compressed)
//          G2*sp*tau (compressed)
//          G2*tp*alpha (compressed)
//          G2*up*beta (compressed)

// G1 and G2 are both arrays of two big.Ints (field elements)
type CompressedG1 big.Int
type CompressedG2 big.Int

type Ppot struct {
	HashOfLastContribution [64]byte
	TauG1                  CompressedG1
	TauG2                  CompressedG2
	AlphaTauG1             CompressedG1
	BetaTauG1              CompressedG1
	PublicKey              PublicKey
}

type PublicKey struct {
	BetaG2    CompressedG2
	G1s       CompressedG1
	G1sTau    CompressedG1
	G1t       CompressedG1
	G1tAlpha  CompressedG1
	G1u       CompressedG1
	G1uBeta   CompressedG1
	G2spTau   CompressedG2
	G2tpAlpha CompressedG2
	G2upBeta  CompressedG2
}

func readPpot() (Ppot, error) {
	return Ppot{}, nil
}
