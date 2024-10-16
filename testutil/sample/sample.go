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

package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

// Account represents a bech32 encoded address and the base64 encoded slice of bytes representing said address.
type Account struct {
	Address   string
	AddressBz []byte
}

// TestAccount returns an Account representing a newly generated PubKey.
func TestAccount() Account {
	pk := ed25519.GenPrivKey().PubKey()
	address := sdk.AccAddress(pk.Address()).String()
	_, bz, _ := bech32.DecodeAndConvert(address)
	return Account{
		Address:   address,
		AddressBz: bz,
	}
}
