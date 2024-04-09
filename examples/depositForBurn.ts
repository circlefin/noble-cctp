/*
 * Copyright (c) 2024, Circle Internet Financial LTD All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
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

import "dotenv/config"
import { DirectSecp256k1HdWallet, Registry, GeneratedType } from "@cosmjs/proto-signing";
import { SigningStargateClient } from "@cosmjs/stargate";
import { MsgDepositForBurn } from "./generated/tx";

export const cctpTypes: ReadonlyArray<[string, GeneratedType]> = [
    ["/circle.cctp.v1.MsgDepositForBurn", MsgDepositForBurn],
];

function createDefaultRegistry(): Registry {
    return new Registry(cctpTypes)
};

const main = async() => {

    const mnemonic = process.env.MNEMONIC ? process.env.MNEMONIC : "";
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        {
            prefix: "noble"
        }
    );

    const [account] = await wallet.getAccounts();

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
            burnToken: "uusdc",
            // If using DepositForBurnWithCaller, add destinationCaller here
        }
    };

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

    console.log(`Burned on Noble: https://mintscan.io/noble-testnet/tx/${result.transactionHash}`);
    console.log(`Minting on Ethereum to https://sepolia.etherscan.io/address/${rawMintRecipient}`);
}

main()