# Noble -> Ethereum (Typescript)

## Instructions

1. Install require packages:

    ```
    npm install
    ```

2. (If needed) Obtain tokens from the faucet: https://faucet.circle.com/

3. Set up `.env` file based on `.env.example`:

    ```
    MNEMONIC="word1 word2..."
    ETH_MINT_RECIPIENT=0x...
    ```

4. Run Noble -> ETH bridging script:

    ```
    npm run burn
    ```

The Noble testnet -> ETH Sepolia CCTP relayer should pick up these messages automatically.
