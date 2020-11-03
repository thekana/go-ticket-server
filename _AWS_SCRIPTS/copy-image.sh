#!/bin/sh
if [ -z "$1" ]
  then
  echo "Require path to key"
  exit 1
fi
if [ -z "$2" ]
  then
  echo "Require user@host"
  exit 1
fi
docker save ticket-reservation | bzip2 | pv | ssh -i $1 $2 'bunzip2 | docker load'