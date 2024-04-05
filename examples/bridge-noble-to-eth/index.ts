import {GeneratedType, Registry} from "@cosmjs/proto-signing";

require("dotenv").config();
const { DirectSecp256k1HdWallet } = require("@cosmjs/proto-signing")
const { SigningStargateClient } = require("@cosmjs/stargate")
const { MsgDepositForBurn } = require("./generated/tx")

export const cctpTypes: ReadonlyArray<[string, GeneratedType]> = [
    ["/circle.cctp.v1.MsgDepositForBurn", MsgDepositForBurn],
];

function createDefaultRegistry(): Registry {
    return new Registry(cctpTypes)
};

const main = async() => {

    const mnemonic = process.env.MNEMONIC;
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        {
            prefix: "noble"
        }
    );

    const [account] = await wallet.getAccounts();

    console.log(account.address)
    const client = await SigningStargateClient.connectWithSigner(
        "https://rpc.testnet.noble.strange.love",
        wallet,
        {
            registry: createDefaultRegistry()
        }
    );

    // Left pad the mint recipient address with 0's to 32 bytes
    const rawMintRecipient = process.env.ETH_MINT_RECIPIENT ? process.env.ETH_MINT_RECIPIENT : "";
    const cleanedMintRecipient = rawMintRecipient.replace(/^0x/, '');
    const zeroesNeeded = 64 - cleanedMintRecipient.length;
    const mintRecipient = '0'.repeat(zeroesNeeded) + cleanedMintRecipient;
    const buffer = Buffer.from(mintRecipient, "hex");
    const mintRecipientBytes = new Uint8Array(buffer);

    const msg = {
        typeUrl: "/circle.cctp.v1.MsgDepositForBurn",
        value: {
            from: account.address,
            amount: "1",
            destinationDomain: 0,
            mintRecipient: mintRecipientBytes,
            burnToken: "uusdc"
        }
    }

    const fee = {
        amount: [
            {
                denom: "uusdc",
                amount: "0",
            },
        ],
        gas: "200000",
    };
    const memo = "";
    const result = await client.signAndBroadcast(
        account.address,
        [msg],
        fee,
        memo
    );

    console.log(`Burned on Noble, tx: https://mintscan.io/noble-testnet/tx/${result.transactionHash}`);
    console.log(`Minting on Ethereum to https://sepolia.etherscan.io/address/${rawMintRecipient}`);
}

main()