package deserializer

///////////////////////////////////////////////////////////////////
///                     PREPARE PHASE 2 PTAU                    ///
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
