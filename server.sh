#!/bin/bash


set -e  # Exit on any error

timestamp=$(date '+%Y-%m-%d %H:%M:%S')
echo "[$timestamp] $1 Starting webserver" >> /var/log/webserver.log

# Define base directory
BASE_DIR="/home/ec2-user/go/src"
API_DIR="$BASE_DIR/stinsondataapi"

mkdir -p $API_DIR

sudo chown -R ec2-user:ec2-user "$API_DIR"

timestamp=$(date '+%Y-%m-%d %H:%M:%S')
echo "[$timestamp] $1 Starting setup" >> /var/log/webserver.log 2>&1

# Clean up existing directory
if [ -d "$API_DIR" ]; then
    rm -rf "$API_DIR"
fi

export GOPRIVATE=github.com/htstinson/business_searcher
git config --global url."git@github.com:".insteadOf "https://github.com/"


# Clone using HTTPS instead of SSH
git clone git@github.com:htstinson/stinsondataapi.git "$API_DIR"


# Create certs directory and copy certificates
mkdir -p "$API_DIR/api/certs"
cp "$BASE_DIR/certs/certificate.crt" "$API_DIR/api/certs/certificate.crt"
cp "$BASE_DIR/certs/private.key" "$API_DIR/api/certs/private.key"

# Change to the API directory
cd "$API_DIR/api/cmd/api"

# Run go mod tidy and start the server
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go build -o webserver .
echo "========================================================================="
sudo /home/ec2-user/go/src/stinsondataapi/api/cmd/api/webserver