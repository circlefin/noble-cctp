package types_test

import (
	"testing"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/stretchr/testify/require"
)

// Message -> bytes -> Message -> bytes
func TestParseMessageHappyPath(t *testing.T) {
	message := &types.Message{
		Version:           1,
		SourceDomain:      2,
		DestinationDomain: 3,
		Nonce:             4,
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
		MessageBody:       []byte("message body"),
	}
	messageBytes, err := message.Bytes()
	require.NoError(t, err)
	parsedMessage, err := new(types.Message).Parse(messageBytes)
	require.NoError(t, err)
	require.Equal(t, message, parsedMessage)
	parsedMessageBytes, err := parsedMessage.Bytes()
	require.NoError(t, err)
	require.Equal(t, messageBytes, parsedMessageBytes)
}

func TestParseIntoMessageWithInvalidInput(t *testing.T) {
	message := types.Message{
		Sender:            []byte("too short"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
	}
	_, err := message.Bytes()
	require.ErrorIs(t, types.ErrParsingMessage, err)

	message = types.Message{
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("too short"),
		DestinationCaller: []byte("destination caller90123456789012"),
	}
	_, err = message.Bytes()
	require.ErrorIs(t, types.ErrParsingMessage, err)

	message = types.Message{
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("too short"),
	}
	_, err = message.Bytes()
	require.ErrorIs(t, types.ErrParsingMessage, err)
}
