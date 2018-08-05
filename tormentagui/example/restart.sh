#!/bin/sh

if [ $# -eq 0 ]
  then
    GO=go
  else
    GO="sysconfcpus -n $1 go"
fi

echo "Changes detected - killing server..."
pkill -f tormentagui
echo "Recompiling and starting..."
# Instead of run, we use build as we need to specify an output name
# because we need to be able to shut the server down by name
$GO build -o=tormentagui *.go
./tormentagui &