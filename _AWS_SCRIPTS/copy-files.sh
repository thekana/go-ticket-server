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

ssh -i $1 $2 'mkdir -p server/docker; mkdir -p server/_dev_server_keys'
scp -i $1 ../docker/docker-compose.yml $2:~/server/docker/
scp -i $1 ../docker/config_docker.yml $2:~/server/docker/
scp -i $1 ../docker/wait-for-it.sh $2:~/server/docker/
scp -i $1 -r ../_dev_server_keys/* $2:~/server/_dev_server_keys
scp -i $1 -r ../nginx $2:~/server