package types

import (
	"encoding/binary"
	"math/big"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Parse parses a byte array into a BurnMessage struct
// https://developers.circle.com/stablecoin/docs/cctp-technical-reference#burnmessage
func (msg *BurnMessage) Parse(bz []byte) (*BurnMessage, error) {
	if len(bz) != BurnMessageLen {
		return nil, sdkerrors.Wrapf(ErrParsingBurnMessage, "burn message must be %d bytes, got %d", BurnMessageLen, len(msg.BurnToken))
	}

	msg.Version = binary.BigEndian.Uint32(bz[BurnMsgVersionIndex:BurnTokenIndex])
	msg.BurnToken = bz[BurnTokenIndex:MintRecipientIndex]
	msg.MintRecipient = bz[MintRecipientIndex:AmountIndex]
	msg.Amount = math.NewIntFromBigInt(new(big.Int).SetBytes(bz[AmountIndex:MsgSenderIndex]))
	msg.MessageSender = bz[MsgSenderIndex:BurnMessageLen]

	return msg, nil
}

// Bytes parses a BurnMessage struct into a byte array
// burn token, mint recipient, and message sender must be 32 bytes
func (msg *BurnMessage) Bytes() ([]byte, error) {
	if len(msg.BurnToken) != BurnTokenLen {
		return nil, sdkerrors.Wrapf(ErrParsingBurnMessage, "burn token must be 32 bytes, got %d", len(msg.BurnToken))
	}
	if len(msg.MintRecipient) != MintRecipientLen {
		return nil, sdkerrors.Wrapf(ErrParsingBurnMessage, "mint recipient must be 32 bytes, got %d", len(msg.MintRecipient))
	}
	if len(msg.MessageSender) != AddressBytesLen {
		return nil, sdkerrors.Wrapf(ErrParsingBurnMessage, "message sender must be 32 bytes, got %d", len(msg.MessageSender))
	}

	result := make([]byte, BurnMessageLen)

	versionBytes := make([]byte, VersionLen)
	binary.BigEndian.PutUint32(versionBytes, msg.Version)

	amountBytes := make([]byte, AmountLen)
	msg.Amount.BigInt().FillBytes(amountBytes)

	copy(result[BurnMsgVersionIndex:BurnTokenIndex], versionBytes)
	copy(result[BurnTokenIndex:MintRecipientIndex], msg.BurnToken)
	copy(result[MintRecipientIndex:AmountIndex], msg.MintRecipient)
	copy(result[AmountIndex:MsgSenderIndex], amountBytes)
	copy(result[MsgSenderIndex:BurnMessageLen], msg.MessageSender)

	return result, nil
}
