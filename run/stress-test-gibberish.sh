#!/bin/bash
# Created by Franco Loyola - For Noroff Final Degree Project
# The goal is to try to overload the slog server, uses the builtin bash port converter

message="<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47 - BOM'su root' failed for lonvick on /dev/pts/: "

port="3333"
ip="127.0.0.1"

nPackets=10000 # Can be the double (if UDP is enabled)

for ((i = 0 ; i <= nPackets ; i++))
do
    echo "$message$i" > /dev/tcp/$ip/$port
    echo "$message$i" > /dev/udp/$ip/$port
    if [[ $i -eq $nPackets ]]
    then
        echo "Sent $i packets on UDP and $i packets on TCP"
    fi
done
