#!/bin/bash
# Created by Franco Loyola - For Noroff Final Degree Project

# Reference vars
startDir=$(pwd)
buildDir="build"
dockerDir="run"
# Script names
buildScript="./build-bins.sh"
#dockerScript="./build-docker-image.sh"
syslogScript="./syslog-gibberish.sh"
slogScript="./slog-gibberish.sh"
fileWriterScript="./test-file-writer.sh"
stressScript="./stress-test-gibberish.sh"
# Replay command
replayCmd="./slog-server-darwin-arm64.bin -password=my-fancy-password -priv-key=keys/storage/slog-storage_rsa -read-file=tmp/slog-logging.slog"

# Show the build process
clear
echo "command: $buildScript"
cd "$buildDir" || echo "Could not change to $buildDir"
$buildScript
echo ""
echo "command: ls -l *.bin"
ls -l *.bin
cd "$startDir" || echo "Could not change to $startDir"
echo "##################################"
read -r

# Client
clear
echo "Start the client in another window so it can be shown"
echo "##################################"
read -r

# Start senders
clear
cd "$dockerDir" || echo "Could not change to $dockerDir"
rm -rf data/slog-server/*
echo "Sender/Traffic generating scripts send two packet a second"
echo "command: run/syslog-gibberish.sh &"
$syslogScript &
echo "command: run/test-file-writer.sh &"
$fileWriterScript &
echo "command: run/slog-gibberish.sh &"
$slogScript &
echo "Packets will be sent, will take a while to complete on purpose... check Grafana :)"
wait
echo "Done sending"
echo "##################################"
read -r

# Replay
clear
cd "$startDir" || echo "Could not change to $startDir"
cd "$buildDir" || echo "Could not change to $buildDir"
echo "Creating tmp folder and move files for the replay showcase"
echo "command: mkdir tmp && cp run/data/slog-server/* tmp/"
rm -rf tmp && mkdir tmp && cp ../run/data/slog-server/* tmp/
echo "command: $replayCmd"
$replayCmd
echo "##################################"
read -r
echo "Simple diff to compare tmp/slog-logging.txt vs slog-reprocessed.txt"
diff "tmp/slog-logging.txt" "slog-reprocessed.txt"
echo "##################################"
read -r

# Modify a line
clear
echo "Modify a line in the .slog file and continue... Should catch it in the replay"
read -r
rm slog-reprocessed.txt
$replayCmd
echo "Diff again to compare tmp/slog-logging.txt vs slog-reprocessed.txt"
diff "tmp/slog-logging.txt" "slog-reprocessed.txt"
echo "Done"
echo "##################################"
read -r

# Stress test
clear
cd "$startDir" || echo "Could not change to $startDir"
cd "$dockerDir" || echo "Could not change to $dockerDir"
echo "Stress test - stop container and use current config"
read -r
echo "command: time run/stress-test-gibberish.sh"
time $stressScript

echo ""
echo "Stop the server, change the config params, run again with bigger buffers"
echo "command: time run/stress-test-gibberish.sh"
read -r
time $stressScript

# Hope all goes OK!
echo "Thanks for the patience! Questions?"
