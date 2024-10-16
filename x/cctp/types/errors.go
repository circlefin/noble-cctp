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

package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/cctp module sentinel errors
var (
	ErrUnauthorized                     = errors.Register(ModuleName, 30, "unauthorized")
	ErrMint                             = errors.Register(ModuleName, 31, "tokens can not be minted")
	ErrBurn                             = errors.Register(ModuleName, 32, "tokens can not be burned")
	ErrAttesterAlreadyFound             = errors.Register(ModuleName, 33, "attester is already present")
	ErrAuthorityNotSet                  = errors.Register(ModuleName, 34, "authority not set")
	ErrMalformedField                   = errors.Register(ModuleName, 35, "field cannot be empty or nil")
	ErrReceiveMessage                   = errors.Register(ModuleName, 36, "err in receive message")
	ErrDisableAttester                  = errors.Register(ModuleName, 37, "err in disable attester")
	ErrUpdateSignatureThreshold         = errors.Register(ModuleName, 38, "err in update signature threshold")
	ErrMinterAllowanceNotFound          = errors.Register(ModuleName, 39, "minter allowance not found")
	ErrTokenPairAlreadyFound            = errors.Register(ModuleName, 40, "token pair already exists")
	ErrTokenPairNotFound                = errors.Register(ModuleName, 41, "token pair not found")
	ErrSendMessage                      = errors.Register(ModuleName, 42, "error in send message")
	ErrSendMessageWithCaller            = errors.Register(ModuleName, 43, "error in send message with caller")
	ErrDepositForBurn                   = errors.Register(ModuleName, 44, "error in deposit for burn")
	ErrInvalidDestinationCaller         = errors.Register(ModuleName, 45, "malformed destination caller")
	ErrSignatureVerification            = errors.Register(ModuleName, 46, "unable to verify signature")
	ErrReplaceMessage                   = errors.Register(ModuleName, 47, "error in replace message")
	ErrDuringPause                      = errors.Register(ModuleName, 48, "error while trying to pause or unpause")
	ErrInvalidAmount                    = errors.Register(ModuleName, 49, "invalid amount")
	ErrNextAvailableNonce               = errors.Register(ModuleName, 50, "error while retrieving next available nonce")
	ErrRemoteTokenMessengerAlreadyFound = errors.Register(ModuleName, 51, "this remote token messenger mapping already exists")
	ErrRemoteTokenMessengerNotFound     = errors.Register(ModuleName, 53, "remote token messenger not found")
	ErrParsingMessage                   = errors.Register(ModuleName, 54, "error while parsing message into bytes")
	ErrParsingBurnMessage               = errors.Register(ModuleName, 55, "error while parsing burn message into bytes")
	ErrInvalidRemoteToken               = errors.Register(ModuleName, 56, "invalid remote token")

	ErrInvalidAddress = errors.Register(ModuleName, 100, "invalid address")
)
