package deserializer

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

type Ppot struct {
}

func readPpot() (Ppot, error) {
	return Ppot{}, nil
}
