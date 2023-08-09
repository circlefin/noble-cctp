package keeper

import (
	"bytes"
	"context"
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) SendMessageWithCaller(goCtx context.Context, msg *types.MsgSendMessageWithCaller) (*types.MsgSendMessageWithCallerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	emptyByteArr := make([]byte, types.DestinationCallerLen)
	if len(msg.DestinationCaller) != types.DestinationCallerLen || bytes.Equal(msg.DestinationCaller, emptyByteArr) {
		return nil, sdkerrors.Wrap(types.ErrSendMessage, "destination caller must be nonzero")
	}

	nonce := k.ReserveAndIncrementNonce(ctx)

	err := k.sendMessage(
		ctx,
		msg.DestinationDomain,
		msg.Recipient,
		msg.DestinationCaller,
		[]byte(msg.From),
		nonce.Nonce,
		msg.MessageBody)

	return &types.MsgSendMessageWithCallerResponse{Nonce: nonce.Nonce}, err
}
