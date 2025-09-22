#!/bin/bash

# Exit on any error
set -e

echo "Building PomoLite..."
go build -o pomo cmd/pomo/main.go

echo "PomoLite built successfully!"

# Check for sudo privileges
if [ "$(id -u)" -eq 0 ]; then
    echo "Installing PomoLite to /usr/local/bin..."
    mv pomo /usr/local/bin/pomo
    echo "PomoLite installed successfully!"
    echo "You can now run 'pomo' from anywhere in your terminal."
else
    echo "You are not running as root. Please move the 'pomo' binary to a directory in your \$PATH."
    echo "For example: sudo mv pomo /usr/local/bin/"
fi
