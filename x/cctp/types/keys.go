package types

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// ModuleName defines the module name
	ModuleName = "cctp"

	// StoreKey defines the primary module store key
	StoreKey = "cctp"

	// RouterKey defines the module's message routing key
	RouterKey = StoreKey

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + StoreKey

	AttesterManagerKey                   = "AttesterManager/value/"
	AuthorityKey                         = "Authority/value/"
	BurningAndMintingPausedKey           = "BurningAndMintingPaused/value/"
	MaxMessageBodySizeKey                = "MaxMessageBodySize/value/"
	NextAvailableNonceKey                = "NextAvailableNonce/value/"
	PauserKey                            = "Pauser/value/"
	SendingAndReceivingMessagesPausedKey = "SendingAndReceivingMessagesPaused/value/"
	SignatureThresholdKey                = "SignatureThreshold/value/"
	TokenControllerKey                   = "TokenController/value/"

	AttesterKeyPrefix            = "Attester/value/"
	PerMessageBurnLimitKeyPrefix = "PerMessageBurnLimit/value/"
	TokenMessengerKeyPrefix      = "TokenMessenger/value/"
	TokenPairKeyPrefix           = "TokenPair/value/"
	UsedNonceKeyPrefix           = "UsedNonce/value/"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// AttesterKey returns the store key to retrieve an Attester from the index fields
func AttesterKey(key []byte) []byte {
	return append(key, []byte("/")...)
}

// PerMessageBurnLimitKey returns the store key to retrieve a PerMessageBurnLimit from the index fields
func PerMessageBurnLimitKey(denom string) []byte {
	return append([]byte(denom), []byte("/")...)
}

// UsedNonceKey returns the store key to retrieve a UsedNonce from the index fields
func UsedNonceKey(nonce uint64, sourceDomain uint32) []byte {
	sourceDomainBz := make([]byte, DomainBytesLen)
	binary.BigEndian.PutUint32(sourceDomainBz, sourceDomain)

	nonceBz := make([]byte, UsedNonceLen)
	binary.BigEndian.PutUint64(nonceBz, nonce)

	result := append(sourceDomainBz, nonceBz...)
	return append(result, []byte("/")...)
}

// TokenPairKey returns the store key to retrieve a TokenPair from the index fields
func TokenPairKey(remoteDomain uint32, remoteToken []byte) []byte {
	remoteDomainBytes := make([]byte, DomainBytesLen)
	binary.BigEndian.PutUint32(remoteDomainBytes, remoteDomain)

	combinedBytes := append(remoteDomainBytes, remoteToken...)
	hashedKey := crypto.Keccak256(combinedBytes)

	return append(hashedKey, []byte("/")...)
}

// TokenMessengerKey returns the store key to retrieve a TokenMessenger from the index fields
func TokenMessengerKey(domain uint32) []byte {
	domainBytes := make([]byte, DomainBytesLen)
	binary.BigEndian.PutUint32(domainBytes, domain)

	return append(domainBytes, []byte("/")...)
}
