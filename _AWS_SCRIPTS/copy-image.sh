#!/bin/sh
if [ $# -eq 0 ]
  then
  echo "Required path to key"
  exit 1
fi
docker save ticket-reservation | bzip2 | pv | ssh -i $1 ec2-user@ec2-52-77-231-100.ap-southeast-1.compute.amazonaws.com 'bunzip2 | docker load'