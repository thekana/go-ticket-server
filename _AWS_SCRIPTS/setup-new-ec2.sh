# For setting up Amazon Linux 2
sudo yum update -y
sudo yum install -y git
sudo amazon-linux-extras install docker
sudo service docker start
sudo usermod -a -G docker ec2-user
sudo curl -L https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# cat setup-new-ec2.sh | ssh -i $KEYPATH ec2-user@ec2-52-77-231-100.ap-southeast-1.compute.amazonaws.com 'bash -s'
