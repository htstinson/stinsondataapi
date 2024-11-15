#!/bin/bash

# Create certs directory if it doesn't exist
mkdir -p certs

# Generate CA key and certificate
openssl req -x509 \
    -sha256 -days 356 \
    -nodes \
    -newkey rsa:2048 \
    -subj "/CN=Local Dev Root CA/C=US" \
    -keyout certs/rootCA.key -out certs/rootCA.crt

# Generate server private key
openssl genrsa -out certs/local.key 2048

# Generate server CSR
openssl req -new \
    -key certs/local.key \
    -out certs/local.csr \
    -subj "/CN=api.local.dev"

# Generate server certificate
openssl x509 -req \
    -sha256 \
    -days 365 \
    -in certs/local.csr \
    -CA certs/rootCA.crt \
    -CAkey certs/rootCA.key \
    -CAcreateserial \
    -out certs/local.crt \
    -extfile <(printf "subjectAltName=DNS:api.local.dev")