#!/bin/bash
# Created by Franco Loyola - For Noroff Final Degree Project
# Script to build both binaries for the slog system, uses relative paths, don't change the file structure!
# If the binary names are changed, update the ../run/build-docker-image.sh and the Dockerfile as well to match it

# Reference vars
startDir=$(pwd)
serverBinName="slog-server"
clientBinName="slog-client"

# Paths, relative to $startDir
serverPath="../code/cmd/slogserver"
clientPath="../code/cmd/slogclient"

# Build combinations
OSlist=("linux" "darwin")
ARCHlist=("386" "amd64" "arm64")

# Show help
echo ""
echo "Builds can be added/removed from this script, check for valid combinations at:"
echo "https://gist.github.com/camabeh/a02e6846e00251e1820c784516c0318f"
echo "If the build fails, check that the combination is valid, not all of them are (see error that pops in build)"
echo ""
echo "This script builds by default for Linux and macOS"

# Remove old bins
echo ""
echo ""
echo "Removing old binaries from this folder..."
rm -f $serverBinName*.bin
rm -f $clientBinName*.bin

# Go vet all
echo ""
echo ""
cd ../
# Static code analysis
echo "Checking with go vet"
if go vet ./...
then
    echo " - go vet passed!"
else
    echo " - Failed to pass go vet"
    exit 1
fi
# Tests
echo "Running tests (if any)"
if go test ./...
then
    echo " - Tests passed!"
else
    echo " - Some tests failed..."
    exit 1
fi

# Server build
echo ""
echo ""
echo "Attempting to build the slog-server"
cd "$startDir" || echo "Could not change to $startDir"
cd "$serverPath" || echo " - Could not find the $serverPath folder..."
for os in "${OSlist[@]}"
do
    for arch in "${ARCHlist[@]}"
    do
        out="$serverBinName-$os-$arch.bin"
        if GOOS=$os GOARCH=$arch go build -o "$out" .
        then
            echo " - '$out built!"
            mv "$out" "$startDir"
        else
            echo " - Failed to build the '$out'"
        fi
    done
done


# Client build
echo ""
echo ""
echo "Attempting to build the slog-client"
cd "$startDir" || echo "Could not change to $startDir"
cd "$clientPath" || echo "Could not find the $clientPath folder..."
for os in "${OSlist[@]}"
do
    for arch in "${ARCHlist[@]}"
    do
        out="$clientBinName-$os-$arch.bin"
        if GOOS=$os GOARCH=$arch go build -o "$out" .
        then
            echo " - '$out built!"
            mv "$out" "$startDir"
        else
            echo " - Failed to build the '$out'"
        fi
    done
done
