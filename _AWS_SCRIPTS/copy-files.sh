#!/bin/sh
if [ -z "$1" ]; then
  echo "Require path to key"
  exit 1
fi
if [ -z "$2" ]; then
  echo "Require user@host"
  exit 1
fi

ssh -i $1 $2 'sudo rm -rf /server; mkdir -p server/docker; mkdir -p server/_dev_server_keys; mkdir -p server/postgres_data'
scp -r -i $1 ../docker $2:~/server
scp -i $1 -r ../_dev_server_keys/* $2:~/server/_dev_server_keys
