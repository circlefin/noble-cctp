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
package types

import (
	"encoding/binary"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Parse parses a byte array into a Message struct
// https://developers.circle.com/stablecoin/docs/cctp-technical-reference#message
func (msg *Message) Parse(bz []byte) (*Message, error) {
	if len(bz) < MessageBodyIndex {
		return nil, sdkerrors.Wrapf(ErrParsingMessage, "cctp message must be at least %d bytes, got %d", MessageBodyIndex, len(bz))
	}

	msg.Version = binary.BigEndian.Uint32(bz[VersionIndex:SourceDomainIndex])
	msg.SourceDomain = binary.BigEndian.Uint32(bz[SourceDomainIndex:DestinationDomainIndex])
	msg.DestinationDomain = binary.BigEndian.Uint32(bz[DestinationDomainIndex:NonceIndex])
	msg.Nonce = binary.BigEndian.Uint64(bz[NonceIndex:SenderIndex])
	msg.Sender = bz[SenderIndex:RecipientIndex]
	msg.Recipient = bz[RecipientIndex:DestinationCallerIndex]
	msg.DestinationCaller = bz[DestinationCallerIndex:MessageBodyIndex]
	msg.MessageBody = bz[MessageBodyIndex:]

	return msg, nil
}

// Bytes parses a Message struct into a byte array
// sender, recipient, destination caller must be 32 bytes
func (msg *Message) Bytes() ([]byte, error) {
	if len(msg.Sender) != AddressBytesLen {
		return nil, sdkerrors.Wrapf(ErrParsingMessage, "sender must be %d bytes, got %d", AddressBytesLen, len(msg.Sender))
	}
	if len(msg.Recipient) != AddressBytesLen {
		return nil, sdkerrors.Wrapf(ErrParsingMessage, "recipient must be %d bytes, got %d", AddressBytesLen, len(msg.Recipient))
	}
	if len(msg.DestinationCaller) != AddressBytesLen {
		return nil, sdkerrors.Wrapf(ErrParsingMessage, "destination caller must be %d bytes, got %d", AddressBytesLen, len(msg.DestinationCaller))
	}

	result := make([]byte, MessageBodyIndex+len(msg.MessageBody))

	versionBytes := make([]byte, VersionLen)
	binary.BigEndian.PutUint32(versionBytes, msg.Version)

	sourceDomainBytes := make([]byte, DomainBytesLen)
	binary.BigEndian.PutUint32(sourceDomainBytes, msg.SourceDomain)

	destinationDomainBytes := make([]byte, DomainBytesLen)
	binary.BigEndian.PutUint32(destinationDomainBytes, msg.DestinationDomain)

	nonceBytes := make([]byte, NonceBytesLen)
	binary.BigEndian.PutUint64(nonceBytes, msg.Nonce)

	copy(result[VersionIndex:SourceDomainIndex], versionBytes)
	copy(result[SourceDomainIndex:DestinationDomainIndex], sourceDomainBytes)
	copy(result[DestinationDomainIndex:NonceIndex], destinationDomainBytes)
	copy(result[NonceIndex:SenderIndex], nonceBytes)
	copy(result[SenderIndex:RecipientIndex], msg.Sender)
	copy(result[RecipientIndex:DestinationCallerIndex], msg.Recipient)
	copy(result[DestinationCallerIndex:MessageBodyIndex], msg.DestinationCaller)
	copy(result[MessageBodyIndex:MessageBodyIndex+len(msg.MessageBody)], msg.MessageBody)

	return result, nil
}
