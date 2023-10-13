# Events

## AttesterEnabled

This event is emitted whenever an attester is enabled. It contains the hex
encoded public key of the attester that was added.

```go
type AttesterEnabled struct {
    Attester string
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgEnableAttester`](./02_messages.md#enableattester)

## AttesterDisabled

This event is emitted whenever an attester is disabled. It contains the hex
encoded public key of the attester that was removed.

```go
type AttesterDisabled struct {
    Attester string
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgDisableAttester`](./02_messages.md#disableattester)

## SignatureThresholdUpdated

This event is emitted whenever the signature threshold is changed. It contains
both the old and new values of the parameter.

```go
type SignatureThresholdUpdated struct {
    OldSignatureThreshold uint64
    NewSignatureThreshold uint64
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUpdateSignatureThreshold`](./02_messages.md#updatesignaturethreshold)

## OwnerUpdated

This event is emitted when the ownership transfer process if finalized. It
contains the old and new owners. See the [`OwnershipTransferStarted`](#ownershiptransferstarted)
event for when this process is started.

```go
type OwnerUpdated struct {
    PreviousOwner string
    NewOwner      string
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgAcceptOwner`](./02_messages.md#acceptowner)

## OwnershipTransferStarted

This event is emitted when the ownership transfer process is started. It
contains the old and potential new owners. See the [`OwnerUpdated`](#ownerupdated)
event for when this process is finalized.

```go
type OwnershipTransferStarted struct {
    PreviousOwner string
    NewOwner      string
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUpdateOwner`](./02_messages.md#updateowner)

## PauserUpdated

This event is emitted when the pauser role is updated. It contains the old and
new pausers.

```go
type PauserUpdated struct {
    PreviousPauser string
    NewPauser      string
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUpdatePauser`](./02_messages.md#updatepauser)

## AttesterManagerUpdated

This event is emitted when the attester manager role is updated. It contains the
old and new attester managers.

```go
type AttesterManagerUpdated struct {
    PreviousAttesterManager string
    NewAttesterManager      string
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUpdateAttestManager`](./02_messages.md#updateattestermanager)

## TokenControllerUpdated

This event is emitted when the token controller role is updated. It contains the
old and new token controllers.

```go
type TokenControllerUpdated struct {
    PreviousTokenController string
    NewTokenController      string
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUpdateTokenController`](./02_messages.md#updatetokencontroller)

## BurningAndMintingPausedEvent

This event is emitted when burning & minting is paused.

```go
type BurningAndMintingPausedEvent struct {}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgPauseBurningAndMinting`](./02_messages.md#pauseburningandminting)

## BurningAndMintingUnpausedEvent

This event is emitted when burning & minting is unpaused.

```go
type BurningAndMintingUnpausedEvent struct {}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUnpauseBurningAndMinting`](./02_messages.md#unpauseburningandminting)

## SendingAndReceivingPausedEvent

This event is emitted when sending & receiving of messages is paused.

```go
type SendingAndReceivingPausedEvent struct {}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgPauseSendingAndReceivingMessages`](./02_messages.md#pausesendingandreceivingmessages)

## SendingAndReceivingUnpausedEvent

This event is emitted when sending & receiving of messages is unpaused.

```go
type SendingAndReceivingUnpausedEvent struct {}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUnpauseSendingAndReceivingMessages`](./02_messages.md#unpausesendingandreceivingmessages)

## DepositForBurn

This event is emitted when a deposit for burn message is sent to a remote domain.

```go
type DepositForBurn struct {
    Nonce                     uint64
    BurnToken                 string
    Amount                    cosmossdk_io_math.Int
    Depositor                 string
    MintRecipient             []byte
    DestinationDomain         uint32
    DestinationTokenMessenger []byte
    DestinationCaller         []byte
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgDepositForBurn`](./02_messages.md#depositforburn)
- [`circle.cctp.v1.MsgDepositForBurnWithCaller`](./02_messages.md#depositforburnwithcaller)
- [`circle.cctp.v1.MsgReplaceDepositForBurn`](./02_messages.md#replacedepositforburn)

## MintAndWithdraw

This event is emitted when a deposit for burn message is received from a remote domain.

```go
type MintAndWithdraw struct {
    MintRecipient string
    Amount        cosmossdk_io_math.Int
    MintToken     string
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgReceiveMessage`](./02_messages.md#receivemessage)

## TokenPairLinked

This event is emitted when a remote token is linked with a local token.

```go
type TokenPairLinked struct {
    LocalToken   string
    RemoteDomain uint32
    RemoteToken  []byte
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgLinkTokenPair`](./02_messages.md#linktokenpair)

## TokenPairUnlinked

This event is emitted when a remote token is unlinked with a local token.

```go
type TokenPairUnlinked struct {
    LocalToken   string
    RemoteDomain uint32
    RemoteToken  []byte
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUnlinkTokenPair`](./02_messages.md#unlinktokenpair)

## MessageSent

This event is emitted when a message is sent to a remote domain.

```go
type MessageSent struct {
    Message []byte
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgDepositForBurn`](./02_messages.md#depositforburn)
- [`circle.cctp.v1.MsgDepositForBurnWithCaller`](./02_messages.md#depositforburnwithcaller)
- [`circle.cctp.v1.MsgReplaceDepositForBurn`](./02_messages.md#replacedepositforburn)
- [`circle.cctp.v1.MsgReplaceMessage`](./02_messages.md#replacemessage)
- [`circle.cctp.v1.MsgSendMessage`](./02_messages.md#sendmessage)
- [`circle.cctp.v1.MsgSendMessageWithCaller`](./02_messages.md#sendmessagewithcaller)

## MessageReceived

This event is emitted when a message is received from a remote domain.

```go
type MessageReceived struct {
    Caller       string
    SourceDomain uint32
    Nonce        uint64
    Sender       []byte
    MessageBody  []byte
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgReceiveMessage`](./02_messages.md#receivemessage)

## MaxMessageBodySizeUpdated

This event is emitted whenever the max message body size is changed. It contains
the new value of the parameter.

```go
type MaxMessageBodySizeUpdated struct {
    NewMaxMessageBodySize uint64
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgUpdateMaxMessageBodySize`](./02_messages.md#updatemaxmessagebodysize)

## RemoteTokenMessengerAdded

This event is emitted when a token messenger is added for a remote domain.

```go
type RemoteTokenMessengerAdded struct {
    Domain               uint32
    RemoteTokenMessenger []byte
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgAddRemoteTokenMessenger`](./02_messages.md#addremotetokenmessenger)

## RemoteTokenMessengerRemoved

This event is emitted when a token messenger is removed for a remote domain.

```go
type RemoteTokenMessengerRemoved struct {
    Domain               uint32
    RemoteTokenMessenger []byte
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.MsgRemoveRemoteTokenMessenger`](./02_messages.md#removeremotetokenmessenger)

## SetBurnLimitPerMessage

This event is emitted whenever the burn limit per message is changed. It
contains the new value of the parameter.

```go
type SetBurnLimitPerMessage struct {
	Token               string
	BurnLimitPerMessage cosmossdk_io_math.Int
}
```

This event is emitted by the following transactions:

- [`circle.cctp.v1.SetBurnLimitPerMessage`](./02_messages.md#setmaxburnamountpermessage)
