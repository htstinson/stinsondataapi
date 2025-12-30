#!/bin/bash
set -e  # Exit on any error

timestamp=$(date '+%Y-%m-%d %H:%M:%S')
echo "[$timestamp] Starting webserver" >> /var/log/webserver.log

# Define base directory
BASE_DIR="/home/ec2-user/go/src"
API_DIR="$BASE_DIR/stinsondataapi"
BUILD_OUTPUT="$API_DIR/webserver"

# Ensure we're not in the directory we're about to delete
cd "$BASE_DIR"

# Clean up existing directory
if [ -d "$API_DIR" ]; then
    echo "Removing existing directory..."
    rm -rf "$API_DIR"
fi

# Configure Git for private repos
export GOPRIVATE=github.com/htstinson/business_searcher
git config --global url."git@github.com:".insteadOf "https://github.com/"

# Clone repository
echo "Cloning repository..."
git clone git@github.com:htstinson/stinsondataapi.git "$API_DIR"

# Create certs directory and copy certificates
mkdir -p "$API_DIR/api/certs"
cp "$BASE_DIR/certs/certificate.crt" "$API_DIR/api/certs/"
cp "$BASE_DIR/certs/private.key" "$API_DIR/api/certs/"

# Change to the API directory
cd "$API_DIR"

# Set Go path
export PATH=$PATH:/usr/local/go/bin

# Run go mod tidy
echo "Tidying dependencies..."
go mod tidy

# Build the binary - output to root of API_DIR
echo "Building webserver..."
go build -o "$BUILD_OUTPUT" ./api/cmd/api

# Verify the binary was created
if [ ! -f "$BUILD_OUTPUT" ]; then
    echo "ERROR: Build failed - webserver binary not created"
    exit 1
fi

# Make it executable
chmod +x "$BUILD_OUTPUT"

echo "========================================================================="
echo "Starting webserver..."
sudo "$BUILD_OUTPUT"