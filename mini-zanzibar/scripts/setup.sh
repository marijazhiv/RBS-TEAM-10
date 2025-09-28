#!/bin/bash

# Mini-Zanzibar Setup Script

echo "Setting up Mini-Zanzibar development environment..."

# Create data directories
echo "Creating data directories..."
mkdir -p data/leveldb
mkdir -p logs

# Install Go dependencies
echo "Installing Go dependencies..."
go mod tidy

# Start Consul in development mode (background)
echo "Starting Consul in development mode..."
if command -v consul &> /dev/null; then
    consul agent -dev -client=0.0.0.0 &
    CONSUL_PID=$!
    echo "Consul started with PID: $CONSUL_PID"
    echo $CONSUL_PID > consul.pid
else
    echo "Warning: Consul not found. Please install Consul or run it separately."
    echo "Download from: https://www.consul.io/downloads"
fi

# Copy environment file
if [ ! -f .env ]; then
    echo "Creating .env file from template..."
    cp .env.example .env
    echo "Please edit .env file with your configuration"
fi

echo "Setup complete!"
echo ""
echo "To start the Mini-Zanzibar server:"
echo "  go run cmd/server/main.go"
echo ""
echo "To stop Consul (if started by this script):"
echo "  kill \$(cat consul.pid) && rm consul.pid"
echo ""
echo "API will be available at: http://localhost:8080"