package keeper

import (
	"bytes"
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) DepositForBurnWithCaller(goCtx context.Context, msg *types.MsgDepositForBurnWithCaller) (*types.MsgDepositForBurnWithCallerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Destination caller must be nonzero. To allow any destination caller, use DepositForBurn().
	emptyByteArr := make([]byte, types.DestinationCallerLen)
	if len(msg.DestinationCaller) == 0 || bytes.Equal(msg.DestinationCaller, emptyByteArr) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDestinationCaller, "invalid destination caller")
	}

	nonce, err := k.depositForBurn(
		ctx,
		msg.From,
		msg.Amount,
		msg.DestinationDomain,
		msg.MintRecipient,
		msg.BurnToken,
		msg.DestinationCaller)

	return &types.MsgDepositForBurnWithCallerResponse{Nonce: nonce}, err
}
