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
 package simulation

import (
	"crypto/elliptic"
	"math/rand"

	"cosmossdk.io/math"

	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func WeightedOperations(_ codec.JSONCodec, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, keeper *keeper.Keeper) simulation.WeightedOperations {
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			1, SimulateAcceptOwner(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateAddRemoteTokenMessenger(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateDepositForBurn(accountKeeper, bankKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateDepositForBurnWithCaller(accountKeeper, bankKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateDisableAttester(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateEnableAttester(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateLinkTokenPair(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulatePausingOfBurningAndMinting(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulatePausingOfSendingAndReceiving(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateReceiveMessage(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateRemoveRemoteTokenMessenger(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateReplaceDepositForBurn(accountKeeper, bankKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateReplaceMessage(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateSendMessage(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateSendMessageWithCaller(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateSetMaxBurnAmountPerMessage(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateUnlinkTokenPair(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateUpdateAttesterManager(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateUpdateMaxMessageBodySize(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateUpdateOwner(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateUpdatePauser(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateUpdateSignatureThreshold(accountKeeper, keeper),
		),
		simulation.NewWeightedOperation(
			1, SimulateUpdateTokenController(accountKeeper, keeper),
		),
	}
}

func SimulateAcceptOwner(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		owner, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetOwner(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		newOwner, _ := simTypes.RandomAcc(r, accounts)

		updateMsg := &types.MsgUpdateOwner{
			From:     owner.Address.String(),
			NewOwner: newOwner.Address.String(),
		}
		acceptMsg := &types.MsgAcceptOwner{From: newOwner.Address.String()}

		updateTx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           updateMsg,
			MsgType:       updateMsg.Type(),
			Context:       ctx,
			SimAccount:    owner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}
		acceptTx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           acceptMsg,
			MsgType:       acceptMsg.Type(),
			Context:       ctx,
			SimAccount:    newOwner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		if operationMsg, futureOperations, err := simulation.GenAndDeliverTx(updateTx, sdk.NewCoins()); err != nil {
			return operationMsg, futureOperations, err
		}
		return simulation.GenAndDeliverTx(acceptTx, sdk.NewCoins())
	}
}

func SimulateAddRemoteTokenMessenger(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		domain := r.Uint32()
		tokenMessenger := make([]byte, 32)
		r.Read(tokenMessenger)

		owner, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetOwner(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgAddRemoteTokenMessenger{
			From:     owner.Address.String(),
			DomainId: domain,
			Address:  tokenMessenger,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    owner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateDepositForBurn(accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		if paused, _ := keeper.GetBurningAndMintingPaused(ctx); paused.Paused {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		caller := accounts[0]

		balance := bankKeeper.GetBalance(ctx, caller.Address, "uusdc")
		keeper.SetPerMessageBurnLimit(ctx, types.PerMessageBurnLimit{
			Denom:  "uusdc",
			Amount: math.NewInt(balance.Amount.Int64()),
		})
		amount := r.Int63n(balance.Amount.Int64())

		domain := r.Uint32()
		mintRecipient := make([]byte, 32)
		r.Read(mintRecipient)
		tokenMessenger := make([]byte, 32)
		r.Read(tokenMessenger)

		keeper.SetRemoteTokenMessenger(ctx, types.RemoteTokenMessenger{
			DomainId: domain,
			Address:  tokenMessenger,
		})

		burn := &types.MsgDepositForBurn{
			From:              caller.Address.String(),
			Amount:            math.NewInt(amount),
			DestinationDomain: domain,
			MintRecipient:     mintRecipient,
			BurnToken:         "uusdc",
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           burn,
			MsgType:       burn.Type(),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateDepositForBurnWithCaller(accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		if paused, _ := keeper.GetBurningAndMintingPaused(ctx); paused.Paused {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		caller := accounts[0]

		balance := bankKeeper.GetBalance(ctx, caller.Address, "uusdc")
		keeper.SetPerMessageBurnLimit(ctx, types.PerMessageBurnLimit{
			Denom:  "uusdc",
			Amount: math.NewInt(balance.Amount.Int64()),
		})
		amount := r.Int63n(balance.Amount.Int64())

		domain := r.Uint32()
		mintRecipient := make([]byte, 32)
		r.Read(mintRecipient)
		destinationCaller := make([]byte, 32)
		r.Read(destinationCaller)
		tokenMessenger := make([]byte, 32)
		r.Read(tokenMessenger)

		keeper.SetRemoteTokenMessenger(ctx, types.RemoteTokenMessenger{
			DomainId: domain,
			Address:  tokenMessenger,
		})

		msg := &types.MsgDepositForBurnWithCaller{
			From:              caller.Address.String(),
			Amount:            math.NewInt(amount),
			DestinationDomain: domain,
			MintRecipient:     mintRecipient,
			BurnToken:         "uusdc",
			DestinationCaller: destinationCaller,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateDisableAttester(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, chainID string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		for i := 0; i < r.Intn(101); i++ {
			enableFunc := SimulateEnableAttester(accountKeeper, keeper)
			if operationMsg, futureOperations, err := enableFunc(r, app, ctx, accounts, chainID); err != nil {
				return operationMsg, futureOperations, err
			}
		}

		attesters := keeper.GetAllAttesters(ctx)
		attester := attesters[r.Intn(len(attesters))]

		attesterManager, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetAttesterManager(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgDisableAttester{
			From:     attesterManager.Address.String(),
			Attester: attester.Attester,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    attesterManager,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateEnableAttester(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		key, _ := crypto.GenerateKey()
		attester := "0x" + common.Bytes2Hex(
			elliptic.Marshal(key.PublicKey, key.PublicKey.X, key.PublicKey.Y), //nolint:staticcheck
		)

		attesterManager, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetAttesterManager(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgEnableAttester{
			From:     attesterManager.Address.String(),
			Attester: attester,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    attesterManager,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateLinkTokenPair(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		domain := r.Uint32()
		token := make([]byte, 32)
		r.Read(token)

		tokenController, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetTokenController(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgLinkTokenPair{
			From:         tokenController.Address.String(),
			RemoteDomain: domain,
			RemoteToken:  token,
			LocalToken:   "uusdc",
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    tokenController,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulatePausingOfBurningAndMinting(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		paused, ok := keeper.GetBurningAndMintingPaused(ctx)
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		pauser, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetPauser(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		var msg sdk.Msg
		if paused.Paused {
			// Then we need to unpause.
			msg = &types.MsgUnpauseBurningAndMinting{
				From: pauser.Address.String(),
			}
		} else {
			// Then we need to pause.
			msg = &types.MsgPauseBurningAndMinting{
				From: pauser.Address.String(),
			}
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       sdk.MsgTypeURL(msg),
			Context:       ctx,
			SimAccount:    pauser,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulatePausingOfSendingAndReceiving(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		paused, ok := keeper.GetSendingAndReceivingMessagesPaused(ctx)
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		pauser, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetPauser(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		var msg sdk.Msg
		if paused.Paused {
			// Then we need to unpause.
			msg = &types.MsgUnpauseSendingAndReceivingMessages{
				From: pauser.Address.String(),
			}
		} else {
			// Then we need to pause.
			msg = &types.MsgPauseSendingAndReceivingMessages{
				From: pauser.Address.String(),
			}
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       sdk.MsgTypeURL(msg),
			Context:       ctx,
			SimAccount:    pauser,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateReceiveMessage(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		if paused, _ := keeper.GetSendingAndReceivingMessagesPaused(ctx); paused.Paused {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		maxBodySize, _ := keeper.GetMaxMessageBodySize(ctx)
		bodySize := r.Intn(int(maxBodySize.Amount))

		sender := make([]byte, 32)
		r.Read(sender)
		recipient := make([]byte, 32)
		r.Read(recipient)
		body := make([]byte, bodySize)
		r.Read(body)

		rawMessage := types.Message{
			Version:           0,
			SourceDomain:      r.Uint32(),
			DestinationDomain: 4,
			Nonce:             0,
			Sender:            sender,
			Recipient:         recipient,
			DestinationCaller: make([]byte, 32),
			MessageBody:       body,
		}
		message, _ := rawMessage.Bytes()

		key, _ := crypto.GenerateKey()
		attester := "0x" + common.Bytes2Hex(
			elliptic.Marshal(key.PublicKey, key.PublicKey.X, key.PublicKey.Y), //nolint:staticcheck
		)

		keeper.SetAttester(ctx, types.Attester{Attester: attester})
		keeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: 1})

		digest := crypto.Keccak256(message)
		attestation, _ := crypto.Sign(digest, key)

		caller, _ := simTypes.RandomAcc(r, accounts)

		msg := &types.MsgReceiveMessage{
			From:        caller.Address.String(),
			Message:     message,
			Attestation: attestation,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       sdk.MsgTypeURL(msg),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateRemoveRemoteTokenMessenger(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, chainID string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		for i := 0; i < r.Intn(101); i++ {
			addFunc := SimulateAddRemoteTokenMessenger(accountKeeper, keeper)
			if operationMsg, futureOperations, err := addFunc(r, app, ctx, accounts, chainID); err != nil {
				return operationMsg, futureOperations, err
			}
		}

		tokenMessengers := keeper.GetRemoteTokenMessengers(ctx)
		domain := tokenMessengers[r.Intn(len(tokenMessengers))].DomainId

		owner, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetOwner(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgRemoveRemoteTokenMessenger{
			From:     owner.Address.String(),
			DomainId: domain,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    owner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateReplaceDepositForBurn(accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, chainID string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		if paused, _ := keeper.GetSendingAndReceivingMessagesPaused(ctx); paused.Paused {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		nonce, _ := keeper.GetNextAvailableNonce(ctx)

		caller := accounts[0]

		balance := bankKeeper.GetBalance(ctx, caller.Address, "uusdc")
		keeper.SetPerMessageBurnLimit(ctx, types.PerMessageBurnLimit{
			Denom:  "uusdc",
			Amount: math.NewInt(balance.Amount.Int64()),
		})
		amount := r.Int63n(balance.Amount.Int64())

		domain := r.Uint32()
		sender := make([]byte, 32)
		copy(sender[12:], caller.Address)
		mintRecipient := make([]byte, 32)
		r.Read(mintRecipient)
		newMintRecipient := make([]byte, 32)
		r.Read(newMintRecipient)
		newDestinationCaller := make([]byte, 32)
		r.Read(newDestinationCaller)
		tokenMessenger := make([]byte, 32)
		r.Read(tokenMessenger)

		keeper.SetRemoteTokenMessenger(ctx, types.RemoteTokenMessenger{
			DomainId: domain,
			Address:  tokenMessenger,
		})

		rawBody := types.BurnMessage{
			Version:       0,
			BurnToken:     crypto.Keccak256([]byte("uusdc")),
			MintRecipient: mintRecipient,
			Amount:        math.NewInt(amount),
			MessageSender: sender,
		}
		body, _ := rawBody.Bytes()

		rawMessage := types.Message{
			Version:           0,
			SourceDomain:      4,
			DestinationDomain: domain,
			Nonce:             nonce.Nonce,
			Sender:            types.PaddedModuleAddress,
			Recipient:         tokenMessenger,
			DestinationCaller: make([]byte, 32),
			MessageBody:       body,
		}
		message, _ := rawMessage.Bytes()

		key, _ := crypto.GenerateKey()
		attester := "0x" + common.Bytes2Hex(
			elliptic.Marshal(key.PublicKey, key.PublicKey.X, key.PublicKey.Y), //nolint:staticcheck
		)

		keeper.SetAttester(ctx, types.Attester{Attester: attester})
		keeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: 1})

		digest := crypto.Keccak256(message)
		attestation, _ := crypto.Sign(digest, key)

		sendMsg := &types.MsgDepositForBurn{
			From:              caller.Address.String(),
			Amount:            math.NewInt(amount),
			DestinationDomain: domain,
			MintRecipient:     mintRecipient,
			BurnToken:         "uusdc",
		}
		replaceMsg := &types.MsgReplaceDepositForBurn{
			From:                 caller.Address.String(),
			OriginalMessage:      message,
			OriginalAttestation:  attestation,
			NewDestinationCaller: newDestinationCaller,
			NewMintRecipient:     newMintRecipient,
		}

		sendTx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           sendMsg,
			MsgType:       sendMsg.Type(),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}
		replaceTx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           replaceMsg,
			MsgType:       replaceMsg.Type(),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		if operationMsg, futureOperations, err := simulation.GenAndDeliverTx(sendTx, sdk.NewCoins()); err != nil {
			return operationMsg, futureOperations, err
		}
		return simulation.GenAndDeliverTx(replaceTx, sdk.NewCoins())
	}
}

func SimulateReplaceMessage(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, chainID string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		if paused, _ := keeper.GetSendingAndReceivingMessagesPaused(ctx); paused.Paused {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		maxBodySize, _ := keeper.GetMaxMessageBodySize(ctx)
		bodySize := r.Intn(int(maxBodySize.Amount))

		nonce, _ := keeper.GetNextAvailableNonce(ctx)

		caller, _ := simTypes.RandomAcc(r, accounts)

		domain := r.Uint32()
		sender := make([]byte, 32)
		copy(sender[12:], caller.Address)
		recipient := make([]byte, 32)
		r.Read(recipient)
		body := make([]byte, bodySize)
		r.Read(body)
		newBody := make([]byte, bodySize)
		r.Read(newBody)
		newDestinationCaller := make([]byte, 32)
		r.Read(newDestinationCaller)

		rawMessage := types.Message{
			Version:           0,
			SourceDomain:      4,
			DestinationDomain: domain,
			Nonce:             nonce.Nonce,
			Sender:            sender,
			Recipient:         recipient,
			DestinationCaller: make([]byte, 32),
			MessageBody:       body,
		}
		message, _ := rawMessage.Bytes()

		key, _ := crypto.GenerateKey()
		attester := "0x" + common.Bytes2Hex(
			elliptic.Marshal(key.PublicKey, key.PublicKey.X, key.PublicKey.Y), //nolint:staticcheck
		)

		keeper.SetAttester(ctx, types.Attester{Attester: attester})
		keeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: 1})

		digest := crypto.Keccak256(message)
		attestation, _ := crypto.Sign(digest, key)

		sendMsg := &types.MsgSendMessage{
			From:              caller.Address.String(),
			DestinationDomain: domain,
			Recipient:         recipient,
			MessageBody:       body,
		}
		replaceMsg := &types.MsgReplaceMessage{
			From:                 caller.Address.String(),
			OriginalMessage:      message,
			OriginalAttestation:  attestation,
			NewMessageBody:       newBody,
			NewDestinationCaller: newDestinationCaller,
		}

		sendTx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           sendMsg,
			MsgType:       sendMsg.Type(),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}
		replaceTx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           replaceMsg,
			MsgType:       replaceMsg.Type(),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		if operationMsg, futureOperations, err := simulation.GenAndDeliverTx(sendTx, sdk.NewCoins()); err != nil {
			return operationMsg, futureOperations, err
		}
		return simulation.GenAndDeliverTx(replaceTx, sdk.NewCoins())
	}
}

func SimulateSendMessage(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, chainID string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		if paused, _ := keeper.GetSendingAndReceivingMessagesPaused(ctx); paused.Paused {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		maxBodySize, _ := keeper.GetMaxMessageBodySize(ctx)
		bodySize := r.Intn(int(maxBodySize.Amount))

		domain := r.Uint32()
		recipient := make([]byte, 32)
		r.Read(recipient)
		body := make([]byte, bodySize)
		r.Read(body)

		caller, _ := simTypes.RandomAcc(r, accounts)

		msg := &types.MsgSendMessage{
			From:              caller.Address.String(),
			DestinationDomain: domain,
			Recipient:         recipient,
			MessageBody:       body,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateSendMessageWithCaller(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, chainID string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		if paused, _ := keeper.GetSendingAndReceivingMessagesPaused(ctx); paused.Paused {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		maxBodySize, _ := keeper.GetMaxMessageBodySize(ctx)
		bodySize := r.Intn(int(maxBodySize.Amount))

		domain := r.Uint32()
		recipient := make([]byte, 32)
		r.Read(recipient)
		body := make([]byte, bodySize)
		r.Read(body)
		destCaller := make([]byte, 32)
		r.Read(destCaller)

		caller, _ := simTypes.RandomAcc(r, accounts)

		msg := &types.MsgSendMessageWithCaller{
			From:              caller.Address.String(),
			DestinationDomain: domain,
			Recipient:         recipient,
			MessageBody:       body,
			DestinationCaller: destCaller,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    caller,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateUnlinkTokenPair(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, chainID string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		linkFunc := SimulateLinkTokenPair(accountKeeper, keeper)
		if operationMsg, futureOperations, err := linkFunc(r, app, ctx, accounts, chainID); err != nil {
			return operationMsg, futureOperations, err
		}

		tokenPairs := keeper.GetAllTokenPairs(ctx)
		tokenPair := tokenPairs[0]

		tokenController, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetTokenController(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgUnlinkTokenPair{
			From:         tokenController.Address.String(),
			RemoteDomain: tokenPair.RemoteDomain,
			RemoteToken:  tokenPair.RemoteToken,
			LocalToken:   tokenPair.LocalToken,
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    tokenController,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateUpdateOwner(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		owner, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetOwner(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		newOwner, _ := simTypes.RandomAcc(r, accounts)

		updateMsg := &types.MsgUpdateOwner{
			From:     owner.Address.String(),
			NewOwner: newOwner.Address.String(),
		}

		updateTx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           updateMsg,
			MsgType:       updateMsg.Type(),
			Context:       ctx,
			SimAccount:    owner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(updateTx, sdk.NewCoins())
	}
}

func SimulateUpdateAttesterManager(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		owner, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetOwner(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		newAttesterManager, _ := simTypes.RandomAcc(r, accounts)

		msg := &types.MsgUpdateAttesterManager{
			From:               owner.Address.String(),
			NewAttesterManager: newAttesterManager.Address.String(),
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    owner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateUpdateTokenController(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		owner, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetOwner(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		newTokenController, _ := simTypes.RandomAcc(r, accounts)

		msg := &types.MsgUpdateTokenController{
			From:               owner.Address.String(),
			NewTokenController: newTokenController.Address.String(),
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    owner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateUpdatePauser(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		owner, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetOwner(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		newPauser, _ := simTypes.RandomAcc(r, accounts)

		msg := &types.MsgUpdatePauser{
			From:      owner.Address.String(),
			NewPauser: newPauser.Address.String(),
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    owner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateUpdateMaxMessageBodySize(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		size := r.Int63n(10_000) // up to 10KB

		owner, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetOwner(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgUpdateMaxMessageBodySize{
			From:        owner.Address.String(),
			MessageSize: uint64(size),
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    owner,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateSetMaxBurnAmountPerMessage(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, _ string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		tokenController, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetTokenController(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgSetMaxBurnAmountPerMessage{
			From:       tokenController.Address.String(),
			LocalToken: "uusdc",
			Amount:     math.NewInt(r.Int63()),
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    tokenController,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}

func SimulateUpdateSignatureThreshold(accountKeeper types.AccountKeeper, keeper *keeper.Keeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simTypes.Account, chainID string) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		for i := 0; i < r.Intn(101); i++ {
			enableFunc := SimulateEnableAttester(accountKeeper, keeper)
			if operationMsg, futureOperations, err := enableFunc(r, app, ctx, accounts, chainID); err != nil {
				return operationMsg, futureOperations, err
			}
		}
		amount := simTypes.RandIntBetween(r, 1, len(keeper.GetAllAttesters(ctx)))

		if signatureThreshold, _ := keeper.GetSignatureThreshold(ctx); signatureThreshold.Amount == uint32(amount) {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		attesterManager, ok := simTypes.FindAccount(accounts, sdk.MustAccAddressFromBech32(keeper.GetAttesterManager(ctx)))
		if !ok {
			return simTypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		msg := &types.MsgUpdateSignatureThreshold{
			From:   attesterManager.Address.String(),
			Amount: uint32(amount),
		}

		tx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         params.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    attesterManager,
			AccountKeeper: accountKeeper,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(tx, sdk.NewCoins())
	}
}
