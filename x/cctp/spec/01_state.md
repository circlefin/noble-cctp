# State

## Roles

### Owner

The owner field is of type `string` and is a Noble address. This field is used
to authenticate all transactions requiring ownership of the CCTP module.

`Key: 0x6f776e6572`

### Pending Owner

The pending owner field is of type `string` and is a Noble address. This field
is used when executing an ownership transfer. The CCTP module utilizes a
two-step transfer flow, where the new owner must first accept ownership before
the transfer is finalized.

`Key: 0x70656e64696e672d6f776e6572`

### Attester Manager

The attester manager field is of type `string` and is a Noble address. This
field is used to authenticate all transactions requiring the attester manager
role, namely:

- [`circle.cctp.v1.MsgDisableAttester`](./02_messages.md#disableattester)
- [`circle.cctp.v1.MsgEnableAttester`](./02_messages.md#enableattester)
- [`circle.cctp.v1.MsgUpdateSignatureThreshold`](./02_messages.md#updatesignaturethreshold)

`Key: 0x61747465737465722d6d616e61676572`

### Pauser

The pauser field is of type `string` and is a Noble address. This field is used
to authenticate all transactions requiring the pauser role, namely:

- [`circle.cctp.v1.MsgPauseBurningAndMinting`](./02_messages.md#pauseburningandminting)
- [`circle.cctp.v1.MsgPauseSendingAndReceivingMessages`](./02_messages.md#pausesendingandreceivingmessages)
- [`circle.cctp.v1.MsgUnpauseBurningAndMinting`](./02_messages.md#unpauseburningandminting)
- [`circle.cctp.v1.MsgUnpauseSendingAndReceivingMessages`](./02_messages.md#unpausesendingandreceivingmessages)

`Key: 0x706175736572`

### Token Controller

The token controller field is of type `string` and is a Noble address. This
field is used to authenticate all transactions requiring the token controller
role, namely:

- [`circle.cctp.v1.MsgLinkTokenPair`](./02_messages.md#linktokenpair)
- [`circle.cctp.v1.MsgSetMaxBurnAmountPerMessage`](./02_messages.md#setmaxburnamountpermessage)
- [`circle.cctp.v1.MsgUnlinkTokenPair`](./02_messages.md#unlinktokenpair)

`Key: 0x746f6b656e2d636f6e74726f6c6c6572`

## Attesters

Attesters are dedicated their own store prefix, which is used to store
individual `Attester` items.

`Key: 0x41747465737465722f76616c75652f`

### Attester

An attester object contains the hex encoded public key of an [Iris API]
attester. These are used to verify message attestations in the
[`circle.cttp.v1.MsgReceiveMessage`](./02_messages.md#receivemessage)
transaction.

```go
type Attester struct {
    Attester string
}
```

`Key: [Attester]/`

## Per Message Burn Limits

Per message burn limits are dedicated their own store prefix, which is used to
store individual `PerMessageBurnLimit` items.

`Key: 0x5065724d6573736167654275726e4c696d69742f76616c75652f`

### `PerMessageBurnLimit`

A per message burn limit object contains the amount of a specific local token is
allowed to be burned per
[`circle.cctp.v1.MsgDepositForBurn`](./02_messages.md#depositforburn) and
[`circle.cctp.v1.MsgDepositForBurnWithCaller`](./02_messages.md#depositforburnwithcaller)
transactions.

```go
type PerMessageBurnLimit struct {
    Denom  string
    Amount cosmossdk_io_math.Int
}
```

`Key: [Denom]/`

## Burning & Minting Paused

The burning & minting paused field is of type `BurningAndMintingPaused`. This
field is used to determine if `BurnMessages` can be sent or received on the
local domain.

```go
type BurningAndMintingPaused struct {
    Paused bool
}
```

`Key: 0x4275726e696e67416e644d696e74696e675061757365642f76616c75652f`

## Sending & Receiving Paused

The sending & receiving paused field is of type
`SendingAndReceivingMessagesPaused`. This field is used to determine if messages
can be sent or received on the local domain.

```go
type SendingAndReceivingMessagesPaused struct {
    Paused bool
}
```

`Key: 0x53656e64696e67416e64526563656976696e674d657373616765735061757365642f76616c75652f`

## Max Message Body Size

The max message body size field is of type `MaxMessageBodySize`. This field is
used when validating messages sent to remote domains.

```go
type MaxMessageBodySize struct {
    Amount uint64
}
```

`Key: 0x4d61784d657373616765426f647953697a652f76616c75652f`

## Next Available Nonce

The next available nonce field is of type `Nonce`. This field is used to
determine the nonce of a message when sending to a remote domain.

```go
type Nonce struct {
    SourceDomain uint32
    Nonce        uint64
}
```

`Key: 0x4e657874417661696c61626c654e6f6e63652f76616c75652f`

## Signature Threshold

The signature threshold field is of type `SignatureThreshold`. This field is
used when authenticating signatures for received messaged.

```go
type SignatureThreshold struct {
    Amount uint32
}
```

`Key: 0x5369676e61747572655468726573686f6c642f76616c75652f`

## Token Pairs

Token pairs are dedicated their own store prefix, which is used to store
individual `TokenPair` items.

`Key: 0x546f6b656e506169722f76616c75652f`

### `TokenPair`

A token pair object contains relevant information surrounding a remote token,
and it's paired local token.

```go
type TokenPair struct {
    RemoteDomain uint32
    RemoteToken  []byte
    LocalToken   string
}
```

`Key: [keccak(RemoteDomain, RemoteToken)]/`

## Used Nonces

Used nonces are dedicated their own store prefix, which is used to store
individual `Nonce` items.

`Key: 0x557365644e6f6e63652f76616c75652f`

### `Nonce`

A nonce object contains relevant information surrounding a used nonce of a
remote domain.

```go
type Nonce struct {
    SourceDomain uint32
    Nonce        uint64
}
```

`Key: [SourceDomain]/[Nonce]/`

## Token Messengers

Token messengers are dedicated their own store prefix, which is used to store
individual `RemoteTokenMessenger` items.

`Key: 0x52656d6f7465546f6b656e4d657373656e6765722f76616c75652f`

### `RemoteTokenMessenger`

A remote token messenger object contains relevant information surrounding a
token messenger sitting on a remote domain.

```go
type RemoteTokenMessenger struct {
    DomainId uint32
    Address  []byte
}
```

`Key: [DomainId]/`

[Iris API]:
https://developers.circle.com/stablecoin/docs#attestation-service-api
