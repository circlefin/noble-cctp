// Copyright 2024 Circle Internet Group, Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package keeper

import (
	"context"

	"cosmossdk.io/errors"
	cctpTypes "github.com/circlefin/noble-cctp/x/cctp/types"
	ftftypes "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MockFiatTokenfactoryKeeper struct{}

var _ cctpTypes.FiatTokenfactoryKeeper = MockFiatTokenfactoryKeeper{}

func (MockFiatTokenfactoryKeeper) Mint(ctx sdk.Context, msg *ftftypes.MsgMint) (*ftftypes.MsgMintResponse, error) {
	return &ftftypes.MsgMintResponse{}, nil
}

func (MockFiatTokenfactoryKeeper) Burn(ctx sdk.Context, msg *ftftypes.MsgBurn) (*ftftypes.MsgBurnResponse, error) {
	return &ftftypes.MsgBurnResponse{}, nil
}

func (MockFiatTokenfactoryKeeper) GetMintingDenom(ctx context.Context) (val ftftypes.MintingDenom) {
	return ftftypes.MintingDenom{Denom: "uusdc"}
}

type MockErrFiatTokenfactoryKeeper struct{}

var _ cctpTypes.FiatTokenfactoryKeeper = MockErrFiatTokenfactoryKeeper{}

func (MockErrFiatTokenfactoryKeeper) Burn(ctx sdk.Context, msg *ftftypes.MsgBurn) (*ftftypes.MsgBurnResponse, error) {
	return &ftftypes.MsgBurnResponse{}, ftftypes.ErrBurn
}

func (k MockErrFiatTokenfactoryKeeper) Mint(ctx sdk.Context, msg *ftftypes.MsgMint) (*ftftypes.MsgMintResponse, error) {
	return nil, errors.Wrap(ftftypes.ErrBurn, "error calling mint")
}

func (k MockErrFiatTokenfactoryKeeper) GetMintingDenom(ctx context.Context) (val ftftypes.MintingDenom) {
	return ftftypes.MintingDenom{Denom: "uusdc"}
}
