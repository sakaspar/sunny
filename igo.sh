#!/bin/bash

# Download the Go binary for ARM64 architecture
echo "Downloading Go binary..."
wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz -P /tmp

# Extract the downloaded tarball to /usr/local
echo "Extracting Go binary..."
sudo tar -C /usr/local -xzf /tmp/go1.22.0.linux-arm64.tar.gz

# Update PATH environment variable
echo "Updating PATH..."
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
echo "Verifying installation..."
go version

echo "Go installation completed."
