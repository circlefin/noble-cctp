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
import { MsgReceiveMessage } from "./generated/tx";

export const cctpTypes: ReadonlyArray<[string, GeneratedType]> = [
    ["/circle.cctp.v1.MsgReceiveMessage", MsgReceiveMessage],
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

    // Convert the message and attestation from hex to bytes
    const messageHex = process.env.MESSAGE_HEX ? process.env.MESSAGE_HEX : "";
    const attestationSignature = process.env.ATTESTATION ? process.env.ATTESTATION : "";

    const messageBytes = new Uint8Array(Buffer.from(messageHex.replace("0x", ""), "hex"));
    const attestationBytes = new Uint8Array(Buffer.from(attestationSignature.replace("0x", ""), "hex"));

    const msg = {
        typeUrl: "/circle.cctp.v1.MsgReceiveMessage",
        value: {
            from: account.address,
            message: messageBytes,
            attestation: attestationBytes,
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

    console.log(`Minted on Noble: https://mintscan.io/noble-testnet/tx/${result.transactionHash}`);
}

main()