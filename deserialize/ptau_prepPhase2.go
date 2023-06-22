package deserializer

///////////////////////////////////////////////////////////////////
///                     PREPARE PHASE 2 PTAU                    ///
///////////////////////////////////////////////////////////////////

// Not needed since gnark does lagrange evaluation on its own
// https://github.com/ConsenSys/gnark/blob/172cc2499244cf9975bdea055aa275ad2230cee9/backend/groth16/bn254/mpcsetup/phase2.go#L102-L106

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
	lagrangeTauG1(12)
		{(2 ** power)*2-1} [
			G1, tau*G1, tau^2 * G1, ....
		]
	lagrangeTauG2(13)
		{2 ** power}[
			G2, tau*G2, tau^2 * G2, ...
		]
	lagrangeAlphaTauG1(14)
		{2 ** power}[
			alpha*G1, alpha*tau*G1, alpha*tau^2*G1,....
		]
	lagrangeBetaTauG1(15)
		{2 ** power} []
			beta*G1, beta*tau*G1, beta*tau^2*G1, ....
		]
*/

// Has 11 sections, has new header, copies sections 2-7 from phase 1
// and writes the lagrange evaluations in sections 12-15 for tauG1,
// tauG2, alphaTauG1 and betaTauG1

type PreparePhase2Ptau struct {
	Header               PtauHeader
	PTauPubKey           PtauPubKey
	LagrangeCoefficients PtauLagrangeCoefficients
}

type PtauLagrangeCoefficients struct {
	LagrangeTauG1      []G1
	LagrangeTauG2      []G2
	LagrangeAlphaTauG1 []G1
	LagrangeBetaTauG1  []G1
}

func ReadPreparePhase2Ptau() (PreparePhase2Ptau, error) {
	return PreparePhase2Ptau{}, nil
}
