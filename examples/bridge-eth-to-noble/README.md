# Ethereum -> Noble (Javascript)

## Instructions

1. Install require packages:

    ```
    npm install
    ```

2. (If needed) Obtain tokens from the faucet: https://faucet.circle.com/

3. Set up `.env` file based on `.env.example`:

    ```
    ETH_TESTNET_RPC=https://sepolia.infura.io/v3/<key>
    ETH_PRIVATE_KEY=0xabc123456...
    NOBLE_ADDRESS=noble...
    ```

4. Run ETH -> Noble bridging script:
    ```
    npm run mint
    ```

The Eth Sepolia -> Noble testnet CCTP relayer should pick these messages up automatically.