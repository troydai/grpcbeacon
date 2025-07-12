#!/bin/bash

# Development Environment Setup Script
# This script installs required tools and sets up the development environment

set -e

echo "🔧 Setting up development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "✅ Go found: $(go version)"

# Get Go binary path
GOPATH=$(go env GOPATH)
GOBIN="${GOPATH}/bin"

# Install buf if not present
if ! command -v buf &> /dev/null; then
    echo "📦 Installing buf..."
    go install github.com/bufbuild/buf/cmd/buf@latest
    echo "✅ buf installed successfully"
else
    echo "✅ buf already installed: $(buf --version)"
fi

# Add Go bin to PATH if not already there
if [[ ":$PATH:" != *":$GOBIN:"* ]]; then
    echo "🔧 Adding Go binary directory to PATH..."
    export PATH="$PATH:$GOBIN"
    
    # Add to shell profile for persistence
    SHELL_RC=""
    if [[ -f "$HOME/.bashrc" ]]; then
        SHELL_RC="$HOME/.bashrc"
    elif [[ -f "$HOME/.zshrc" ]]; then
        SHELL_RC="$HOME/.zshrc"
    fi
    
    if [[ -n "$SHELL_RC" ]]; then
        if ! grep -q "export PATH.*$GOBIN" "$SHELL_RC"; then
            echo "export PATH=\"\$PATH:$GOBIN\"" >> "$SHELL_RC"
            echo "✅ Added Go binary path to $SHELL_RC"
        fi
    fi
else
    echo "✅ Go binary directory already in PATH"
fi

# Verify tools are working
echo "🧪 Verifying installation..."

if command -v buf &> /dev/null; then
    echo "✅ buf is working: $(buf --version)"
else
    echo "❌ buf installation failed"
    exit 1
fi

echo "🎉 Development environment setup complete!"
echo ""
echo "💡 You can now run:"
echo "   make gen    # Generate protocol buffer code"
echo "   make bin    # Build the server binary"
echo "   make test   # Run tests"
echo ""
echo "🔄 Note: If this is your first time running the setup, you may need to:"
echo "   source ~/.bashrc  # or restart your terminal"