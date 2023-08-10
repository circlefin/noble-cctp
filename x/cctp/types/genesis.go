package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                            DefaultParams(),
		Authority:                         nil,
		AttesterList:                      []Attester{},
		PerMessageBurnLimitList:           []PerMessageBurnLimit{},
		BurningAndMintingPaused:           &BurningAndMintingPaused{Paused: false},
		SendingAndReceivingMessagesPaused: &SendingAndReceivingMessagesPaused{Paused: false},
		MaxMessageBodySize:                nil,
		NextAvailableNonce:                nil,
		SignatureThreshold:                nil,
		TokenPairList:                     []TokenPair{},
		UsedNoncesList:                    []Nonce{},
		TokenMessengerList:                []TokenMessenger{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.  Stateful checks are performed in InitGenesis
func (gs GenesisState) Validate() error {

	if gs.Authority != nil {
		_, err := sdk.AccAddressFromBech32(gs.Authority.Address)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
		}
	}

	// Check for duplicated index in attesters
	attesterIndexMap := make(map[string]struct{})
	for _, elem := range gs.AttesterList {
		index := string(AttesterKey([]byte(elem.Attester)))
		if _, ok := attesterIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for attesters")
		}
		attesterIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in per message burn limit
	perMessageBurnLimitIndexMap := make(map[string]struct{})
	for _, elem := range gs.PerMessageBurnLimitList {
		index := string(PerMessageBurnLimitKey(elem.Denom))
		if _, ok := perMessageBurnLimitIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for per message burn limits")
		}
		perMessageBurnLimitIndexMap[index] = struct{}{}
	}

	if gs.BurningAndMintingPaused == nil {
		return fmt.Errorf("BurningAndMintingPaused cannot be nil")
	}

	if gs.SendingAndReceivingMessagesPaused == nil {
		return fmt.Errorf("SendingAndReceivingMessagesPaused cannot be nil")
	}

	// Check for duplicated index in token pairs
	tokenPairIndexMap := make(map[string]struct{})
	for _, elem := range gs.TokenPairList {
		index := string(TokenPairKey(elem.RemoteDomain, elem.RemoteToken))
		if _, ok := attesterIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for token pairs")
		}
		tokenPairIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in the used nonce list
	usedNonceIndexMap := make(map[string]struct{})
	for _, elem := range gs.UsedNoncesList {
		index := string(UsedNonceKey(elem.Nonce, elem.SourceDomain))
		if _, ok := usedNonceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for used nonces")
		}
		usedNonceIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in token messengers
	tokenMessengerIndexMap := make(map[string]struct{})
	for _, elem := range gs.TokenMessengerList {
		index := string(TokenMessengerKey(elem.DomainId))
		if _, ok := tokenMessengerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for token messengers")
		}
		tokenMessengerIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
