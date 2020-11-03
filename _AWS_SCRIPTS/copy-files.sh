#!/bin/sh
if [ $# -eq 0 ]
  then
  echo "Required path to key"
  exit 1
fi
ssh -i $1 ec2-user@ec2-52-77-231-100.ap-southeast-1.compute.amazonaws.com 'mkdir -p server/docker; mkdir -p server/_dev_server_keys'
scp -i $1 ../docker/docker-compose.yml ec2-user@ec2-52-77-231-100.ap-southeast-1.compute.amazonaws.com:~/server/docker/
scp -i $1 ../docker/config_docker.yml ec2-user@ec2-52-77-231-100.ap-southeast-1.compute.amazonaws.com:~/server/docker/
scp -i $1 ../docker/wait-for-it.sh ec2-user@ec2-52-77-231-100.ap-southeast-1.compute.amazonaws.com:~/server/docker/
scp -i $1 -r ../_dev_server_keys/* ec2-user@ec2-52-77-231-100.ap-southeast-1.compute.amazonaws.com:~/server/_dev_server_keys
scp -i $1 -r ../nginx ec2-user@ec2-52-77-231-100.ap-southeast-1.compute.amazonaws.com:~/server