#!/bin/bash
# Created by Franco Loyola - For Noroff Final Degree Project
# Script to build the Docker image to run the slog server within a container
# If the binary names are changed, update the ../build/build-bins.sh and the Dockerfile as well to match it
# If the Docker iamge name changes, update the docker-compose.yml as well to match it

# Reference vars
startDir=$(pwd)
arch="arm64"
os="linux"
serverBinName="slog-server-$os-$arch.bin"
serverConfFile="slog-server.conf"
buildDir="../build"
keysFolder="../build/keys"

# Help
echo ""
echo "The script defaults the build to the ARM build"
echo "Change this to match your architecture locally at the start of this script."
echo "If the values are updated, please update the Dockerfile to match the new ones."

# First we need to build our binary
cd "$buildDir" || echo "Could not change to $buildDir"
bash build-bins.sh
cd "$startDir" || echo "Could not change to $startDir"

# Move the built binary and it's config file for the Dockerfile
cp "$buildDir"/"$serverBinName" .
cp "$buildDir"/"$serverConfFile" .
cp -r "$keysFolder" .

# Build the Docker image
docker build -t "slog-server" .

if [ $? -ne 0 ]; then
    echo "Something went wrong with docker build..."
    exit 1
fi

# Remove the copied file once is not required
rm -rf ./*.bin ./*.conf "keys"

echo ""
echo "All built!"
echo "You can now run the Docker image with your desired password/keys by mounting a local volume into the container"
echo "IMPORTANT : If the data is to be preserved, mount a local volume -v ..."
echo "docker run -v /path/to/your/keys/folder:/app/keys/ -v /path/to/your/longterm/storage:/app/slog-data slog-server -password some-secure-password"
echo ""
echo "NOTE: The paths for the container/binary for the keys and data can be changed in the config file"
