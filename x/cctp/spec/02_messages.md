# Messages

## AcceptOwner

`MsgAcceptOwner` 

Broadcasts a transaction that finalizes a transfer of ownership.
It accepts the `newOwner` and completes the two stage process of transferring an owner. 

This message accepts no arguments.

Requires:
   - [MsgUpdateOwner](#updateowner) must be broadcasted first. This sets a pending `newOwner`
   - Message must be sent from the [`Pending Owner`](./01_state.md#pending-owner) account

State changes: 
   - [`owner`](./01_state.md#owner)

Events emitted:
   - [`OwnerUpdated`](./03_events.md#ownerupdated)


## AddRemoteTokenMessenger

`MsgAddRemoteTokenMessenger`

Broadcast a transaction that adds a remote token messenger for a provided domain.

Arguments:
   - `domain-id`
   - `address` - address of [`RemoteTokenMessenger`](./01_state.md#remotetokenmessenger) to add

Requires:
   - Must be sent from [`Owner`](./01_state.md#owner) account

State Changes:
   - [`RemoteTokenMessenger`](./01_state.md#remotetokenmessenger)

Events emitted:
   - [`RemoteTokenMessengerAdded`](./03_events.md#remotetokenmessengeradded)


## DepositForBurn

`MsgDepositForBurn`

Broadcast a transaction that deposits for burn to a provided domain.

Arguments:
   - `Amount` - The burn amount 
   - `DestinationDomain` - Domain of destination chain
   - `MintRecipient` - address receiving minted tokens on destination chain as a 32 length byte array
   - `BurnToken` - The burn token address on source domain

Requires:
   - `Amount` must be positive
   - `Amount` must be <= [`PerMessageBurnLimit`](./01_state.md#permessageburnlimit)
   - `MintRecipient` must not be blank
   - `BurnToken` must align with a denom that can be minted from the `fiattokenfactory`
   - [`BurningAndMintingPaused`](./01_state.md#burning--minting-paused) must be false

State Changes:
   - (*indirect*)[`NextAvailableNonce`](./01_state.md#next-available-nonce) (Changed by: [`SendMessage`](#sendmessage))

Events emitted:
   - [`DepositForBurn`](./03_events.md#depositforburn)
   - (*indirect*) [`MessageSent`](./03_events.md#messagesent) (Emitted by: [`SendMessage`](#sendmessage))

If a Destination Caller is not included, this message calls: [`SendMessage`](#sendmessage)

If a Destination Caller is included, this message calls: [`SendMessageWithCaller`](#sendmessagewithcaller). See [DepositForBurnWithCaller](#depositforburnwithcaller)


## DepositForBurnWithCaller

`MsgDepositForBurnWithCaller`

Broadcast a transaction that deposits for burn with caller to a provided domain.

This message wraps [`MsgDepositForBurn`](#depositforburn). It adds one extra argument:

Arguments:
   - `DestinationCaller` - authorized caller as 32 length byte array of receiveMessage() on destination domain

Requires:
   - `DestinationCaller` must not be blank


## DisableAttester

`MsgDisableAttester`

Broadcast a transaction that disables a provided attester.

Arguments:
   - `attester` - address of [attester](./01_state.md#attester) to disable.

Requires:
   - Message must be sent from the [`Attester Manager`](./01_state.md#attester-manager) account
   - You cannot remove an attester if there is currently only one attester. `len(storedAttesters) > 0`
   - Number of current attesters must be greater than the signature threshold; disallow removing public key 
   if it causes the n in m/n multisig to fall below m (threshold # of signers)

State changes:
   - [`attester`](./01_state.md#attester)

Events emitted:
   - [`AttesterDisabled`](./03_events.md#attesterdisabled)


## EnableAttester

`MsgEnableAttester`

Broadcast a transaction that enables a provided attester.

Arguments:
   - `attester` - address of [attester](./01_state.md#attester) to enable.

Requires:
   - Message must be sent from the [`Attester Manager`](./01_state.md#attester-manager) account

State changes:
   - [`attester`](./01_state.md#attester)

Events emitted:
   - [`AttesterEnabled`](./03_events.md#attesterenabled)

## LinkTokenPair

`MsgLinkTokenPair`

Broadcast a transaction that links a token pair for a provided domain. This is used for minting/burning.
It maps a remote token on a remote domain to a local token.

Arguments:
   - `LocalToken` - Denom of local token in uunits
   - `RemoteToken` - The remote token address
   - `RemoteDomain` - The domain where the message originated from.

Requires:
   - Message must be sent from the [`Token Controller`](./01_state.md#token-controller) account

State changes:
   - [`TokenPair`](./01_state.md#tokenpair)

Events emitted:
   - [`TokenPairLinked`](./03_events.md#tokenpairlinked)

## PauseBurningAndMinting

`MsgPauseBurningAndMinting`

Broadcast a transaction that pauses burning & minting.

This message accepts no arguments.

Requires:
   - Message must be sent from the [`Pauser`](./01_state.md#pauser) account

State changes:
   - [`BurningAndMintingPaused`](./01_state.md#burning--minting-paused)

Events emitted:
   - [`BurningAndMintingPausedEvent`](./03_events.md#burningandmintingpausedevent)


## PauseSendingAndReceivingMessages

`MsgPauseSendingAndReceivingMessages`

Broadcast a transaction that pauses sending & receiving messages.

This message accepts no arguments.

Requires:
   - Message must be sent from the [`Pauser`](./01_state.md#pauser) account

State changes:
   - [`SendingAndReceivingMessagesPaused`](./01_state.md#sending--receiving-paused)

Events emitted:
   - [`SendingAndReceivingPausedEvent`](./03_events.md#sendingandreceivingpausedevent)


## ReceiveMessage

`MsgReceiveMessage`

Broadcast a transaction that receives a provided message from another domain. After validation, it performs a mint.

Arguments:
   - `message` (https://developers.circle.com/stablecoin/docs/cctp-technical-reference#message)
   - `attestation` - Concatenated 65-byte signature(s) of `message`, in increasing order
      of the attester address recovered from signatures. See [Valid Attestation](#valid-attestation)

Requires:
   - [`SendingAndReceivingMessagesPaused`](./01_state.md#sending--receiving-paused) must be false
   - [`BurningAndMintingPaused`](./01_state.md#burning--minting-paused) must be false
   - if `message` includes destination caller, then message must be sent from the same address as destination caller
   - `message.destinationDomain` must be `4` (Noble chain's domain-ID)
   - `message.version` must be equal to Noble chain's message version (`0`)
   - `message.nonce` must not be a [used nonce](./01_state.md#used-nonces)
   - if `message.messageBody` is a valid [`burnMessage`](https://developers.circle.com/stablecoin/docs/cctp-technical-reference#burnmessage), then:
      - `burnMessage.version` must be equal to the hard coded `MessageBodyVersion` (`0`)
      - `burnMessage.burnToken` and `message.sourceDomain` must be a valid [`token pair`](./01_state.md#tokenpair)


State changes:
    - [`nonce`](./01_state.md#used-nonces) - sets a used nonce


Events emitted: 
   - [`MintAndWithdraw`](./03_events.md#mintandwithdraw)
   - [`MessageReceived`](./03_events.md#messagereceived)

### Valid Attestation

1. Length of `_attestation` == 65 (signature length) * signatureThreshold
2. Addresses recovered from attestation must be in increasing order.
      For example, if signature A is signed by address 0x1..., and signature B
 		is signed by address 0x2..., attestation must be passed as AB.
3. No duplicate signers
4. All signers must be enabled attesters

## RemoveRemoteTokenMessenger

`MsgRemoveRemoteTokenMessenger`

Broadcast a transaction that removes the remote token messenger of a provided domain.

Arguments:
   - `domainId` 

Requires:
   - Must be sent from [`Owner`](./01_state.md#owner) account

State Changes:
   - [Token Messengers](./01_state.md#remotetokenmessenger)

Events emitted:
   - [`RemoteTokenMessengerRemoved`](./03_events.md#remotetokenmessengerremoved)

## ReplaceDepositForBurn

`MsgReplaceDepositForBurn`

Broadcast a transaction that replaces a deposit for burn message. Replace the mint recipient and/or
destination caller.

Allows the sender of a previous BurnMessage (created by depositForBurn or depositForBurnWithCaller)
to send a new BurnMessage to replace the original. The new BurnMessage will reuse the amount and 
burn token of the original without requiring a new deposit.

Arguments:
   - `OriginalMessage`- original message bytes to replace
   - `OriginalAttestation`- attestation bytes of `OriginalMessage`
   - `NewDestinationCaller` - the new destination caller, which may be the
     same as the original destination caller, a new destination caller, or an empty
     destination caller, indicating that any destination caller is valid.
   - `NewMintRecipient` - the new mint recipient. May be the same as the
     original mint recipient, or different.

Requires:
   - [`BurningAndMintingPaused`](./01_state.md#burning--minting-paused) must be false
   - Must be sent from the same account as the original message sender

State Changes:
   - (*indirect*)[`NextAvailableNonce`](./01_state.md#next-available-nonce) (Changed by calling [`MsgReplaceMessage`](#replacemessage)
   which calls: [`SendMessage`](#sendmessage))

Events emitted: 
   - [`DepositForBurn`](./03_events.md#depositforburn)
   - (*indirect*) [`MessageSent`](./03_events.md#messagesent) (Emitted by calling [`MsgReplaceMessage`](#replacemessage)
   which calls: [`SendMessage`](#sendmessage))

 This message calls: [`MsgReplaceMessage`](#replacemessage)


## ReplaceMessage

`MsgReplaceMessage`

Broadcast a transaction that replaces a provided message. Replace the message body and/or destination caller.

Arguments:
   - `OriginalMessage` - original message bytes to replace
   - `OriginalAttestation` - attestation bytes of `OriginalMessage`
   - `NewMessageBody` - new message body of replaced message
   - `NewDestinationCaller` - the new destination caller, which may be the
     same as the original destination caller, a new destination caller, or an empty
     destination caller, indicating that any destination caller is valid.

Requires:
   - [`SendingAndReceivingMessagesPaused`](./01_state.md#sending--receiving-paused) must be false
   - The attestation signatures of the original message must still be valid. Changing attesters or the signature threshold can render all previous messages irreplaceable
   - Must be sent from the same account as the original message sender
   - The `OriginalMessage` `sourceDomain` must be equal to Noble Chain's source domain (`4`)

State Changes:
   - (*indirect*)[`NextAvailableNonce`](./01_state.md#next-available-nonce) (Changed by: [`SendMessage`](#sendmessage))
   
Events emitted: 
   - (*indirect*) [`MessageSent`](./03_events.md#messagesent) (Emitted by: [`SendMessage`](#sendmessage))

 This message calls: [`SendMessage`](#sendmessage)


## SendMessage

`MsgSendMessage`

Broadcast a transaction that sends a message to a provided domain.

Arguments:
   - `DestinationDomain` - Domain of destination chain
   - `Recipient` - Address of message recipient on destination chain
   - `MessageBody` - Raw bytes content of message

Requires:
   - [`SendingAndReceivingMessagesPaused`](./01_state.md#sending--receiving-paused) must be false
   - if [`MaxMessageBodySize`](./01_state.md#max-message-body-size) is set, check that the `MessageBody` 
   is not greater than the[`MaxMessageBodySize`](./01_state.md#max-message-body-size).
   - `Recipient` cannot be blank

State Changes:
   - [`NextAvailableNonce`](./01_state.md#next-available-nonce)

Events emitted:
   - [`MessageSent`](./03_events.md#messagesent)

## SendMessageWithCaller

`MsgSendMessageWithCaller`

Broadcast a transaction that sends a message with a caller to a provided domain.

Specifying a Destination caller requires that only the specified caller can call receiveMessage() on destination domain.

This message wraps [`SendMessage`](#sendmessage). It adds one extra argument:

Arguments:
   - `DestinationCaller` - caller on the destination domain, as 32 length byte array

Requires:
   - `DestinationCaller` cannot be blank

State Changes:
   - [`NextAvailableNonce`](./01_state.md#next-available-nonce)

## SetMaxBurnAmountPerMessage

`MsgSetMaxBurnAmountPerMessage`

Broadcast a transaction that updates the max burn amount per message for a provided token.

Arguments:
   - `LocalToken` - Denom of local token in uunits
   - `Amount`

Requires:
   - Must be sent from [`Token Controller`](./01_state.md#token-controller) account

State Changes:
   - [`PerMessageBurnLimit`](./01_state.md#permessageburnlimit)

Events emitted:
   - [`SetBurnLimitPerMessage`](./03_events.md#setburnlimitpermessage)

## UnlinkTokenPair

`MsgUnlinkTokenPair`

Broadcast a transaction that unlinks a token pair for a provided domain.

Arguments:
   - `RemoteDomain` - The domain tied to remote token
   - `RemoteToken`- The remote token address
   - `LocalToken` - Denom of local token in uunits

Requires:
   - Must be sent from [`Token Controller`](./01_state.md#token-controller) account
   - The `RemoteDomain` and `RemoteToken` must correlate to a valid [tokenPair](./01_state.md#tokenpair)

   State Changes:
   - [`TokenPair`](./01_state.md#tokenpair)

   Events emitted:
      - [`TokenPairUnlinked`](./03_events.md#tokenpairunlinked)

## UpdateAttesterManager

`MsgUpdateAttesterManager`

Broadcast a transaction that updates the attester manager to the provided address.

Arguments:
   - `NewAttesterManager` - address of the new attester manager

Requires:
   - Must be sent from [`Owner`](./01_state.md#owner) account

State Changes:
   - [`attesterManager`](./01_state.md#attester-manager)

Events Emitted:
   [`AttesterManagerUpdated`](./03_events.md#attestermanagerupdated)

## UpdateMaxMessageBodySize

`MsgUpdateMaxMessageBodySize`

Broadcast a transaction that updates the max message body size to the provided size.

Arguments:
   - `MessageSize` - new max message body size

Requires:
   - Must be sent from [`Owner`](./01_state.md#owner) account

State Changes:
   - [`MaxMessageBodySize`](./01_state.md#max-message-body-size)

Events emitted:
   - [`MaxMessageBodySizeUpdated`](./03_events.md#maxmessagebodysizeupdated)

## UpdateOwner

`MsgUpdateOwner`

Broadcast a transaction that initiates a transfer of ownership to the provided address.

Arguments:
   - `newOwner` - noble address to set as the [`Pending Owner`](./01_state.md#pending-owner)

Requires:
   - Must be sent from [`Owner`](./01_state.md#owner) account

State Changes:
   - [`pendingOwner`](./01_state.md#pending-owner)

Events emitted:
   - [`OwnershipTransferStarted`](./03_events.md#ownershiptransferstarted)


## UpdatePauser

`MsgUpdatePauser`

Broadcast a transaction that updates the pauser to the provided address.

Arguments:
   - `NewPauser`- address of the new pauser

Requires:
   - Must be sent from [`Owner`](./01_state.md#owner) account

State Changes:
   - [`pauser`](./01_state.md#pauser)

Events emitted:
   - [`PauserUpdated`](./03_events.md#pauserupdated)

## UpdateSignatureThreshold

`MsgUpdateSignatureThreshold`

Broadcast a transaction that updates the signature threshold to the provided amount.

Arguments:
   - `Amount`

Requires:
   - Must be sent from [`Attester Manager`](./01_state.md#attester-manager) account
   - The new signature threshold cannot be greater than the number of [attesters](./01_state.md#attesters)

State Changes:
   - [`SignatureThreshold`](./01_state.md#signature-threshold)

Events emitted:
   - [`SignatureThresholdUpdated`](./03_events.md#signaturethresholdupdated)

## UpdateTokenController

`MsgUpdateTokenController`

Broadcast a transaction that updates the token controller to the provided address.

Arguments:
   - [`NewTokenController`](./01_state.md#token-controller) - address of the new token controller

Requires:
   - Must be sent from [`Owner`](./01_state.md#owner) account

State Changes:
   - [`TokenController`](./01_state.md#token-controller)

Events emitted:
   - [`TokenControllerUpdated`](./03_events.md#tokencontrollerupdated)

## UnpauseBurningAndMinting

`MsgUnpauseBurningAndMinting`

Broadcast a transaction that unpauses burning & minting.

This message accepts no arguments.

Requires:
   - Message must be sent from the [`Pauser`](./01_state.md#pauser) account

State changes:
   - [`BurningAndMintingPaused`](./01_state.md#burning--minting-paused)

Events emitted:
   - [`BurningAndMintingUnpausedEvent`](./03_events.md#burningandmintingunpausedevent)

## UnpauseSendingAndReceivingMessages

`MsgUnpauseSendingAndReceivingMessages`

Broadcast a transaction that unpauses sending & receiving messages.

This message accepts no arguments.

Requires:
   - Message must be sent from the [`Pauser`](./01_state.md#pauser) account

State changes:
   - [`SendingAndReceivingMessagesPaused`](./01_state.md#sending--receiving-paused)

Events emitted:
   - [`SendingAndReceivingUnpausedEvent`](./03_events.md#sendingandreceivingunpausedevent)
