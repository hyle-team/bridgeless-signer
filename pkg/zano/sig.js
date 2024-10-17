const axios = require('axios');
const { ethers } = require("ethers");
const { exit } = require('process');
const fs = require('fs');

const noderpc =""

const firstWallet  = "http://localhost:10500/json_rpc"
const secondWallet = "http://localhost:10505/json_rpc"
const secondKey    = ""
const firstKey     = ""

let hashes = []


/// Define an async function that takes method name and parameters
async function callJsonRpc(requestData, link) {
    try {
        const response = await axios.post(link, requestData, {
            headers: {
                'Content-Type': 'application/json',
            },
        });

        // Return the result from the JSON-RPC response
        return response.data;
    } catch (error) {
        console.error('Error:', error);
        throw error; // Optionally rethrow the error to handle it in the calling function
    }
}

function reverseHexBytes(hexString) {
    // Ensure the hex string length is even
    if (hexString.length % 2 !== 0) {
        throw new Error("Invalid hex string length");
    }

    // Split the hex string into chunks of 2 characters (1 byte)
    const bytes = hexString.match(/.{1,2}/g);

    // Reverse the array of bytes and join them back into a string
    const reversedHex = bytes.reverse().join('');

    return reversedHex;
}

async function deploy_asset()
{
    try {
        //Generated Private Key: 0x17a938099954cee510d7fc9eb2366f0762b093d9be547acabf8be85f774ef154
        //Generated Address: 0x0886bA9F5b117D2A3C1ce18106F2Ce759f5D34C8

        const loadedWallet = new ethers.Wallet(secondKey);
        console.log("Loaded Address:", loadedWallet.address);
        console.log("Public key:", loadedWallet.signingKey.compressedPublicKey);
        const owner_eth_pub_key = loadedWallet.signingKey.compressedPublicKey.substring(2);
        console.log("Generated Public key HEX:", owner_eth_pub_key);
        const jsonObject = {
            id: 0,
            jsonrpc: "2.0",
            method: "deploy_asset",
            params: {
                asset_descriptor: {
                    //current_supply: 1000000000000000,
                    decimal_point: 12,
                    full_name: "TEST Zano wrapped ABC",
                    hidden_supply: false,
                    meta_info: "Stable and private",
                    owner: "",
                    ticker: "TZABC",
                    total_max_supply: 1000000000000000000,
                    owner_eth_pub_key: owner_eth_pub_key
                },
                destinations: [
                    {
                        address: "ZxCuASh1nm3PzJThzA5fJf6BpnpBVZCj8iyKYEtnMsVBXzjxrvuh5X9TmWzDxezSPKjJLzAscgtgFWNJRKMT2WZL16fiGvfdm",
                        amount: 1000000000000000,
                        asset_id: ""
                    }
                ],
                do_not_split_destinations: false
            }
        };

        const res = await callJsonRpc(jsonObject, secondWallet);
        console.log("deploy_asset response: " + JSON.stringify(res, null, 2));
        /*
        deploy_asset response:
            {
              "id": 0,
              "jsonrpc": "2.0",
              "result": {
                "new_asset_id": "7d51ecaad2e3458e0d62b146f33079c6ea307841b09a44b777e0c01eb11b98bf",
                "tx_id": "73ff52bf4d85153f2b25033dd76e9e92e63214ed983682182e6e2b2ce0ecf46c"
              }
            }
        */
    }

    catch (error) {
        console.error('Error occurred:', error);
    }
}

async function emmit_asset(asset_id)
{
    try {

        //Generated Private Key: 0x17a938099954cee510d7fc9eb2366f0762b093d9be547acabf8be85f774ef154
        //Generated Address: 0x0886bA9F5b117D2A3C1ce18106F2Ce759f5D34C8
        // asset_id 7d51ecaad2e3458e0d62b146f33079c6ea307841b09a44b777e0c01eb11b98bf

        //var use_pregenerated_files = false;

        const loadedWallet = new ethers.Wallet(secondKey);

        console.log("Loaded Address:", loadedWallet.address);
        console.log("Public key:", loadedWallet.signingKey.compressedPublicKey);
        const owner_eth_pub_key = loadedWallet.signingKey.compressedPublicKey.substring(2);
        console.log("Generated Public key HEX:", owner_eth_pub_key);

        //////////////////////////////////////////////////////////////////////
        //////////////////////////////////////////////////////////////////////
        //this part is performed on coordinator node:

        var res_emmit;
        //if(!use_pregenerated_files)
        //{
        const requestDataEmit = {
            id: 0,
            jsonrpc: "2.0",
            method: "emit_asset",
            params: {
                asset_id: asset_id,
                destinations: [{
                    address: "ZxDphM9gFU359BXfg2BsPi4xrfapivmTi1c1pvogvD3dbAdha4iCosCWup8YkyitrvdAH15Cin65C2AFpA3AF6cJ2amvcNF7w",
                    amount: 100000000000,
                    asset_id: ""
                }],
                do_not_split_destinations: false
            }
        };

        res_emmit = await callJsonRpc(requestDataEmit, secondWallet);
        //fs.writeFileSync('emmit_response.json', JSON.stringify(res_emmit, null, 2));
        console.log("emmit_response response: " + JSON.stringify(res_emmit, null, 2));
        //}else
        //{
        //  const data = fs.readFileSync('emmit_response.json', 'utf8');
        //  res_emmit = JSON.parse(data);
        //}


        //////////////////////////////////////////////////////////////////////
        //////////////////////////////////////////////////////////////////////
        //this part is performed on validator node:


        var res_decrypt;
        //if(!use_pregenerated_files)
        //{
        const requestDataDecrypt = {
            id: 0,
            jsonrpc: "2.0",
            method: "decrypt_tx_details",
            params: {
                outputs_addresses: res_emmit.result.data_for_external_signing.outputs_addresses,
                tx_blob: res_emmit.result.data_for_external_signing.unsigned_tx,
                tx_id: "",
                tx_secret_key: res_emmit.result.data_for_external_signing.tx_secret_key
            }
        };

        res_decrypt = await callJsonRpc(requestDataDecrypt, noderpc); //request to daemon
        //fs.writeFileSync('decrypt_response.json', JSON.stringify(res_decrypt, null, 2));
        console.log("decrypt_response : " + JSON.stringify(res_decrypt, null, 2));

        //TODO: response holds all information about what this transaction actually transfer and to what addresses

        //}else
        //{
        //  const data = fs.readFileSync('decrypt_response.json', 'utf8');
        //  res_decrypt = JSON.parse(data);
        //}


        //////////////////////////////////////////////////////////////////////
        //////////////////////////////////////////////////////////////////////
        //this part is performed with TSS scheme:
        const bytesToSign = ethers.getBytes('0x' + res_decrypt.result.verified_tx_id);
        const signature = loadedWallet.signingKey.sign(bytesToSign).serialized;

        //////////////////////////////////////////////////////////////////////
        //////////////////////////////////////////////////////////////////////
        //this part is performed on coordinator node with given signature:

        const signature_without_0x = signature.substring(2);
        console.log("Generated signature: " + signature_without_0x);
        // Strip the last byte (recovery parameter) to get 64 bytes
        const strippedSignature = signature_without_0x.slice(0, -2); // Remove the last byte (2 hex chars)

        const requestSendSigned = {
            id: 0,
            jsonrpc: "2.0",
            method: "send_ext_signed_asset_tx",
            params: {
                eth_sig: strippedSignature,
                expected_tx_id: res_decrypt.result.verified_tx_id,
                finalized_tx: res_emmit.result.data_for_external_signing.finalized_tx,
                unlock_transfers_on_fail: false,
                unsigned_tx: res_emmit.result.data_for_external_signing.unsigned_tx
            }
        }

        const res_sign = await callJsonRpc(requestSendSigned, secondWallet);
        //fs.writeFileSync('sign_response.json', JSON.stringify(res_sign, null, 2));
        console.log("sign_response response: " + JSON.stringify(res_sign, null, 2));
        hashes.push(res_decrypt.result.verified_tx_id)
    }
    catch (error) {
        console.error('Error occurred:', error);
    }
}

async function test_sign() {
    const id = "70dd03634d73880375109e0e6a57fb2769f83562d229646b32ab8a8362d932bf"

    const loadedWallet = new ethers.Wallet(secondKey);
    const bytesToSign = ethers.getBytes('0x' + id);
    const signature = loadedWallet.signingKey.sign(bytesToSign).serialized;

    //////////////////////////////////////////////////////////////////////
    //////////////////////////////////////////////////////////////////////
    //this part is performed on coordinator node with given signature:

    const signature_without_0x = signature.substring(2);
    console.log("Generated signature: " + signature_without_0x);

    const final_res = "0ffc63e5113ee4d8da4262f1d9f4a0c6d6340e7042070c5e72bb748430ce17b456dcf27bf375be49a0071f9234778123e6dbb4f8ffd65f2989b5593864a2c0b41c"
}

async function whitelist_add(id) {
    const request_send = {
        id: 0,
        jsonrpc: "2.0",
        method: "assets_whitelist_add",
        params: {
            "asset_id": id
        }
    }

    const resp = await callJsonRpc(request_send, firstWallet)
    console.log("whitelist resp: "+JSON.stringify(resp, null, 2))
}

async function balance() {
    const request_send = {
        id: 0,
        jsonrpc: "2.0",
        method: "getbalance",
        params: {

        }
    }

    const resp = await callJsonRpc(request_send, firstWallet)
    console.log("whitelist resp: "+JSON.stringify(resp, null, 2))
}

async function main()
{
    try {


         //await deploy_asset();
        // TODO: wait for 10 confirmations
        // //wait for 10 confirmations
        //
        // //asset id 7d51ecaad2e3458e0d62b146f33079c6ea307841b09a44b777e0c01eb11b98bf
        // */
        await emmit_asset("693271afac4c26971fe3ebe5d87853c625795966960ff05588ae42450fb96074");
        await emmit_asset("693271afac4c26971fe3ebe5d87853c625795966960ff05588ae42450fb96074");
        //await test_sign()
        //await whitelist_add("cab92cb5338d7b9f533c404c884cadb3ba579601074a6216e28f0d4da13e2c14")
        //await balance()
        console.log(hashes)




    } catch (error) {
        console.error('Error occurred:', error);
    }

}


main();