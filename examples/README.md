# Noble <-> Ethereum Bridging (Typescript)

## DepositForBurn instructions

1. Install required packages:

    ```
    npm install
    ```

2. (If needed) Obtain tokens from the faucet: https://faucet.circle.com/

3. Set up a `.env` file based on `.env.example`, filling in the `MNEMONIC` and `ETH_MINT_RECIPIENT` fields:

    ```
    MNEMONIC="word1 word2..."
    ETH_MINT_RECIPIENT=0x...
    ```

4. Run the depositForBurn script:

    ```
    npm run depositForBurn
    ```

The Noble testnet -> ETH Sepolia CCTP relayer should pick up these messages automatically. To avoid these being automatically picked up, all references to `MsgDepositForBurn` can be changed to `MsgDepositForBurnWithCaller` and a `destinationCaller` field should be added to `msg.value` below line 70.

## ReceiveMessage instructions

1. Install required packages:

    ```
    npm install
    ```

2. Initiate a `DepositForBurnWithCaller` from ETH to Noble. If a regular `DepositForBurn` call is made, the relayer will automatically receive the message on Noble.

3. Fetch the attestation and message from Iris at https://iris-api-sandbox.circle.com/messages/{sourceDomain}/{txHash}.

4. Set up a `.env` file based on `.env.example`, filling in the `MNEMONIC`, `ATTESTATION`, and `MESSAGE_HEX` fields:

    ```
    MNEMONIC="word1 word2..."
    ATTESTATION=0x...
    MESSAGE_HEX=0x
    ```

5. Run the receiveMessage script:

    ```
    npm run receiveMessage
    ```