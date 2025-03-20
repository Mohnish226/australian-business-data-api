#!/bin/bash

BUILD_DIR=build

if ! [ -x "$(command -v go)" ]; then
    echo "Error: go is not installed, Please follow the instructions at https://go.dev/doc/install" >&2
    exit 1
fi

if [ ! -d $BUILD_DIR ]; then
    mkdir $BUILD_DIR
fi

go mod tidy

if [ "$1" == "--help" ]; then
    echo "Usage: ./build.sh [option]"
    echo "Options:"
    echo "  --all: Build for all platforms (amd64)"
    echo "  --linux: Build for Linux (amd64)"
    echo "  --windows: Build for Windows (amd64)"
    echo "  --mac: Build for macOS (amd64)"
    exit 0
fi

function linux {
    GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/australian-business-data-api-linux ./cmd/australian-business-data-api
    echo "Linux build complete"
}

function windows {
    GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/australian-business-data-api-windows.exe ./cmd/australian-business-data-api
    echo "Windows build complete"
}

function mac {
    GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/australian-business-data-api-mac ./cmd/australian-business-data-api
    echo "macOS build complete"
}

if [ "$1" == "--all" ]; then
    echo "Building for all platforms"
    linux
    windows
    mac
    echo "All builds complete in the build directory"
    exit 0
fi

if [ "$1" == "--linux" ]; then
    linux
    exit 0
fi

if [ "$1" == "--windows" ]; then
    windows
    exit 0
fi

if [ "$1" == "--mac" ]; then
    mac
    exit 0
fi

echo "Building for current platform"
go build -o $BUILD_DIR/australian-business-data-api ./cmd/australian-business-data-api
echo "Build Completed and in folder /build"
