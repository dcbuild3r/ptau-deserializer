package deserializer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeserializer(t *testing.T) {
	input_path := "semaphore_16.zkey"

	assert := require.New(t)

	zkey, err := ReadZkey(input_path)

	if err != nil {
		assert.NoError(err)
	}

	fmt.Println(zkey.Header.ProtocolID)

	// protocolID should be 1 (Groth16)
	assert.Equal(GROTH_16_PROTOCOL_ID, zkey.Header.ProtocolID)
}
