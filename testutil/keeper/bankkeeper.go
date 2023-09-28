/*
 * Copyright (c) 2023, © Circle Internet Financial, LTD.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type MockBankKeeper struct{}

func (MockBankKeeper) SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return nil
}

func (MockBankKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (MockBankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (MockBankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (MockBankKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return nil
}

func (MockBankKeeper) GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool) {
	return banktypes.Metadata{}, true
}

func (MockBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return sdk.Coin{
		Denom:  "uusdc",
		Amount: sdk.Int{},
	}
}

//

type ErrBankKeeper struct{}

func (ErrBankKeeper) SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return nil
}

func (ErrBankKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (ErrBankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (ErrBankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (ErrBankKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return errors.New("intentional error")
}

func (ErrBankKeeper) GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool) {
	return banktypes.Metadata{}, true
}

func (ErrBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return sdk.Coin{
		Denom:  "uusdc",
		Amount: sdk.Int{},
	}
}
