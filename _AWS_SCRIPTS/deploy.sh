#!/bin/sh
if [ -z "$1" ]; then
  echo "Require path to key"
  exit 1
fi
if [ -z "$2" ]; then
  echo "Require user@host"
  exit 1
fi
if ! [ -z "$3" ]; then
  echo "BUILDING DOCKER IMAGE"
  sh ../docker/build.sh
  echo "BUILD COMPLETE"
else
  echo "SKIP BUILDING"
fi

# Docker custom aliases
ssh -i $1 $2 'sudo rm -rf server; docker stop $(docker ps -aq); docker rm $(docker ps -aq);'
echo "COPYING IMAGE TO AWS"
sh copy-image.sh $1 $2
echo "COPYING DEPENDENCIES TO AWS"
sh copy-files.sh $1 $2
echo "DOCKER-COMPOSE"
ssh -i $1 $2 'cd server/docker; docker-compose up -d'
#ssh -i $1 $2 'cd server/docker; docker-compose up --scale api=3 -d'
