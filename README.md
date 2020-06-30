# Overview

SilentDragonLite (SDL) is a fork of silentdragonlite lightwalletd, which is a fork of [lightwalletd](https://github.com/adityapk00/lightwalletd) from the ECC. 

It is a backend service that provides a bandwidth-efficient interface to the Hush blockchain for the [SilentDragonLite cli](https://github.com/MyHush/silentdragonlite-light-cli).

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

#### 1. Run a Hush node.
Start a `hushd` with the following options:
```
server=1
rpcuser=user
rpcpassword=password
rpcbind=127.0.0.1
txindex=1
```

You might need to run with `-reindex` the first time if you are enabling the `txindex` or `insightexplorer` options for the first time. The reindex might take a while.

#### 2. Get a TLS certificate
##### "Let's Encrypt" certificate using NGINX as a reverse proxy
If you running a public-facing server, the easiest way to obtain a certificate is to use a NGINX reverse proxy and get a Let's Encrypt certificate. [Instructions are here](https://www.nginx.com/blog/using-free-ssltls-certificates-from-lets-encrypt-with-nginx/)

Create a new section for the NGINX reverse proxy:
```
server {
    listen 443 ssl http2;
 
 
    ssl_certificate     ssl/cert.pem; # From certbot
    ssl_certificate_key ssl/key.pem;  # From certbot
    
    location / {
        # Replace localhost:9067 with the address and port of your gRPC server if using a custom port
        grpc_pass grpc://localhost:9067;
    }
}
```

#### 3. Run the frontend:
You can run the gRPC server with or without TLS, depending on how you configured step 2. If you are using NGINX as a reverse proxy and are letting NGINX handle the TLS authentication, then run the frontend with `-no-tls`

```
go run cmd/server/main.go -bind-addr 127.0.0.1:9067 -conf-file ~/.komodo/HUSH3/HUSH3.conf -no-tls
```

If you have a certificate that you want to use (either self signed, or from a certificate authority), pass the certificate to the frontend:

```
go run cmd/server/main.go -bind-addr 127.0.0.1:443 -conf-file ~/.komodo/HUSH3/HUSH3.conf  -tls-cert /etc/letsencrypt/live/YOURWEBSITE/fullchain.pem -tls-key /etc/letsencrypt/live/YOURWEBSITE/privkey.pem
```

You should start seeing the frontend ingest and cache the zcash blocks after ~15 seconds. 

#### 4. Point the `silentdragonlite-cli` to this server
Connect to your server!
```
./silentdragonlite-cli -server https://lite.myhush.org
```
