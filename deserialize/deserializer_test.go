package deserializer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Format
// ======
// Header(1)
//      Prover Type 1 Groth
// HeaderGroth(2)
//      n8q
//      q
//      n8r
//      r
//      NVars
//      NPub
//      DomainSize  (multiple of 2
//      alpha1
//      beta1
//      delta1
//      beta2
//      gamma2
//      delta2

func TestDeserializerPtau(t *testing.T) {
	input_path := "powersOfTau28_hez_final_08.ptau"

	assert := require.New(t)

	ptau, err := ReadPtau(input_path)

	if err != nil {
		assert.NoError(err)
	}

	fmt.Printf("Size of the primes in bytes: %v \n", ptau.Header.n8)
}

func TestDeserializerZkey(t *testing.T) {
	input_path := "semaphore_16.zkey"

	assert := require.New(t)

	zkey, err := ReadZkey(input_path)

	if err != nil {
		assert.NoError(err)
	}

	fmt.Printf("ProtocolID for Groth16: %v \n", zkey.ZkeyHeader.ProtocolID)

	// protocolID should be 1 (Groth16)
	assert.Equal(GROTH_16_PROTOCOL_ID, zkey.ZkeyHeader.ProtocolID)

	fmt.Printf("n8q is: %v \n", zkey.protocolHeader.n8q)

	fmt.Printf("q is: %v \n", zkey.protocolHeader.q.String())

	fmt.Printf("n8r is: %v \n", zkey.protocolHeader.n8r)

	fmt.Printf("r is: %v \n", zkey.protocolHeader.r.String())

	fmt.Printf("nVars is: %v \n", zkey.protocolHeader.nVars)

	fmt.Printf("nPublic is: %v \n", zkey.protocolHeader.nPublic)

	fmt.Printf("domainSize is: %v \n", zkey.protocolHeader.domainSize)

	fmt.Printf("power is: %v \n", zkey.protocolHeader.power)
}
