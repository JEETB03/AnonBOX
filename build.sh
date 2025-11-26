#!/bin/bash

echo "🏴‍☠️  AnonBOX Build Script 🏴‍☠️"
echo "================================="

# Check for Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed! Please install Go from https://go.dev/dl/"
    exit 1
fi
echo "✅ Go is installed."

# Check for GCC
if ! command -v gcc &> /dev/null; then
    echo "⚠️  GCC is not installed! GUI build might fail."
    echo "   Please install build-essential (Linux) or Xcode Command Line Tools (macOS)."
    BUILD_GUI=0
else
    echo "✅ GCC is installed."
    BUILD_GUI=1
fi

echo ""
echo "📦 Building CLI..."
go build -o anonbox-cli ./cmd/cli
if [ $? -ne 0 ]; then
    echo "❌ CLI Build Failed!"
    exit 1
fi
echo "✅ CLI Built: ./anonbox-cli"

if [ $BUILD_GUI -eq 1 ]; then
    echo ""
    echo "🎨 Building GUI..."
    go build -o anonbox-gui ./cmd/gui
    if [ $? -ne 0 ]; then
        echo "❌ GUI Build Failed!"
    else
        echo "✅ GUI Built: ./anonbox-gui"
    fi
fi

echo ""
echo "🎉 Build Complete!"
