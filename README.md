# Overview

SilentDragonLite (SDL) is a fork of Zecwallet lightwalletd, which is a fork of [lightwalletd](https://github.com/adityapk00/lightwalletd) from the ECC. 

It is a backend service that provides a bandwidth-efficient interface to the Hush blockchain for the [SilentDragonLite cli](https://github.com/MyHush/zecwallet-light-cli).

## Changes from upstream lightwalletd
This version of lightwalletd extends lightwalletd and:

* Adds support for HUSH
* Adds support for transparent addresses
* Adds several new RPC calls for lightclients
* Lots of perf improvements
  * Replaces SQLite with in-memory cache for Compact Blocks
  * Replace local Txstore, delegating Tx lookups to Zcashd
  * Remove the need for a separate ingestor

## Running your own SDL lightwalletd

#### 0. First, install [Go >= 1.11](https://golang.org/dl/#stable).

#### 1. Generate a TLS self-signed certificate
```
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```
Answer the certificate questions to generate the self-signed certificate

#### 2. You need to run a full node with the following options in HUSH3.conf
```
server=1
rpcuser=user
rpcpassword=password
rpcbind=127.0.0.1
txindex=1
```

#### 3. Run the frontend:
You'll need to use the certificate generated from step 1
```
go run cmd/server/main.go -bind-addr 127.0.0.1:9069 -conf-file ~/.komodo/HUSH3/HUSH3.conf  -tls-cert cert.pem -tls-key key.pem
```

#### 4. Point the `silentdragon-cli` to this server
```
./silentdragon-cli --server https://127.0.0.1:9069 --dangerous
```
