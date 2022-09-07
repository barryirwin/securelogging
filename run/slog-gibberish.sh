#!/bin/bash
# Created by Franco Loyola - For Noroff Final Degree Project
# For the demo.sh

message="Invalid packet arriving to slog-server port"
port="3333"
ip="127.0.0.1"

nPackets=200

for ((i = 0 ; i <= nPackets ; i++))
do
    echo "$message$i" > /dev/tcp/$ip/$port
    echo "$message$i" > /dev/udp/$ip/$port
    sleep 0.1
    if [[ $i -eq $nPackets ]]
    then
        echo "Sent $i packets on UDP and $i packets on TCP"
    fi
done
