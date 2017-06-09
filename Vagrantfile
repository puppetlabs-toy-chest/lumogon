# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.

DOCKER_VERSION="1.12.6-0~ubuntu-xenial"
#DOCKER_VERSION="1.13.1-0~ubuntu-xenial"
#DOCKER_VERSION="17.03.0~ce-0~ubuntu-xenial"
#DOCKER_VERSION="17.04.0~ce-0~ubuntu-xenial"
#DOCKER_VERSION="17.05.0~ce-0~ubuntu-xenial"

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.synced_folder ".", "/home/ubuntu/go/src/github.com/puppetlabs/lumogon"
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
  end

  config.vm.provision "shell", inline: <<-SHELL
    export GOPATH=/home/ubuntu/go
    export PATH=$PATH:$GOPATH/bin
    sudo chown ubuntu:ubuntu -R $GOPATH
    sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
    echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" | tee /etc/apt/sources.list.d/docker.list
    sudo add-apt-repository ppa:longsleep/golang-backports
    sudo apt-get update
    sudo apt-get install -y --allow-downgrades golang-go docker-engine=#{DOCKER_VERSION}
    docker rmi puppet/lumogon
    cd /home/ubuntu/go/src/github.com/puppetlabs/lumogon && rm -rf vendor && make all

    docker version
    sudo docker run --rm -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon --disable-analytics version
    sudo docker run -d nginx

    sudo docker run --rm -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon --disable-analytics scan
  SHELL
end
