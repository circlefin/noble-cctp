package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// BurnMessage -> bytes -> BurnMessage -> bytes
func TestParseIntoBurnMessageHappyPath(t *testing.T) {
	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     crypto.Keccak256([]byte("usdc")),
		MintRecipient: []byte("recipient01234567890123456789012"),
		Amount:        math.NewInt(345678),
		MessageSender: []byte("message-sender567890123456789012"),
	}
	burnMessageBytes, err := burnMessage.Bytes()
	require.Nil(t, err)
	parsedBurnMessage, err := new(types.BurnMessage).Parse(burnMessageBytes)
	require.Nil(t, err)

	require.Equal(t, burnMessage.Version, parsedBurnMessage.Version)
	require.Equal(t, burnMessage.BurnToken, parsedBurnMessage.BurnToken)
	require.Equal(t, burnMessage.MintRecipient, parsedBurnMessage.MintRecipient)
	require.Equal(t, burnMessage.Amount, parsedBurnMessage.Amount)
	require.Equal(t, burnMessage.MessageSender, parsedBurnMessage.MessageSender)

	parsedBurnMessageBytes, err := parsedBurnMessage.Bytes()
	require.Nil(t, err)
	require.Equal(t, burnMessageBytes, parsedBurnMessageBytes)
}

func TestParseIntoBurnMessageWithInvalidInput(t *testing.T) {
	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     []byte("too short"),
		MintRecipient: []byte("recipient01234567890123456789012"),
		Amount:        math.NewInt(345678),
		MessageSender: []byte("message-sender567890123456789012"),
	}
	_, err := burnMessage.Bytes()
	require.ErrorIs(t, types.ErrParsingBurnMessage, err)

	burnMessage = types.BurnMessage{
		Version:       1,
		BurnToken:     crypto.Keccak256([]byte("usdc")),
		MintRecipient: []byte("too short"),
		Amount:        math.NewInt(345678),
		MessageSender: []byte("message-sender567890123456789012"),
	}
	_, err = burnMessage.Bytes()
	require.ErrorIs(t, types.ErrParsingBurnMessage, err)

	burnMessage = types.BurnMessage{
		Version:       1,
		BurnToken:     crypto.Keccak256([]byte("usdc")),
		MintRecipient: []byte("recipient01234567890123456789012"),
		Amount:        math.NewInt(345678),
		MessageSender: []byte("too short"),
	}
	_, err = burnMessage.Bytes()
	require.ErrorIs(t, types.ErrParsingBurnMessage, err)
}
