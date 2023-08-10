package types

import (
	"encoding/hex"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RemoteTokenPadded returns the remote token as a byte array, padded to 32 bytes
func RemoteTokenPadded(remoteTokenHex string) ([]byte, error) {
	remoteToken, err := hex.DecodeString(strings.TrimPrefix(remoteTokenHex, "0x"))
	if err != nil {
		return nil, sdkerrors.Wrap(ErrInvalidRemoteToken, "must be hex string")
	}

	if len(remoteToken) > BurnTokenLen {
		return nil, sdkerrors.Wrapf(ErrInvalidRemoteToken, "must be less than %d bytes", BurnTokenLen)
	}

	remoteTokenPadded := make([]byte, BurnTokenLen)
	for i := 0; i < BurnTokenLen-len(remoteToken); i++ {
		remoteTokenPadded[i] = 0
	}
	copy(remoteTokenPadded[BurnTokenLen-len(remoteToken):], remoteToken)

	return remoteTokenPadded, nil
}
