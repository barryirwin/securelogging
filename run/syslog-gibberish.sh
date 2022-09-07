#!/bin/bash
# Created by Franco Loyola - For Noroff Final Degree Project
# For the demo.sh

message="some important Syslog message: "

port="514"
ip="127.0.0.1"

nPackets=200 # Can be the double

for ((i = 0 ; i <= nPackets ; i++))
do
    echo "$message$i TCP" > /dev/tcp/$ip/$port
    echo "$message$i UDP" > /dev/udp/$ip/$port
    sleep 0.1
    if [[ $i -eq $nPackets ]]
    then
        echo "Sent $i packets on UDP and $i packets on TCP"
    fi
done
