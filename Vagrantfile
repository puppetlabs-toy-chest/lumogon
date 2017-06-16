require 'yaml'

MACHINE_CONFIG = <<-YAML
---
ubuntu:
  xenial:
    box: ubuntu/xenial64
    docker:
      1.12.6: 1.12.6-0~ubuntu-xenial
      1.13.1: 1.13.1-0~ubuntu-xenial
      17.03: 17.03.0~ce-0~ubuntu-xenial
      17.04: 17.04.0~ce-0~ubuntu-xenial
      17.05: 17.05.0~ce-0~ubuntu-xenial
centos:
  6:
    box: puppetlabs/centos-6.6-64-puppet
    docker:
      1.7.0: 1.7.0-1.el6
      1.7.1: 1.7.1-1.el6
  7:
    box: puppetlabs/centos-7.2-64-puppet
    docker:
      1.7.1: 1.7.1-1.el7.centos
      1.8.3: 1.8.3-1.el7.centos
      1.9.1: 1.1.1-1.el7.centos
      1.10.3: 1.10.3-1.el7.centos
      1.11.2: 1.11.2-1.el7.centos
      1.12.6: 1.12.6-1.el7.centos
      1.13.1: 1.13.1-1.el7.centos
      17.03: 17.03.1.ce-1.el7.centos
      17.04: 17.04.0.ce-1.el7.centos
      17.05: 17.05.0.ce-1.el7.centos
YAML

Vagrant.configure("2") do |config|
  machines = YAML.load(MACHINE_CONFIG)

  machines['ubuntu'].each do |distro, matrix|
    matrix['docker'].each do |version, package|
      config.vm.define "ubuntu-#{distro}-#{version}" do |m|
        m.vm.box = matrix['box']
        m.vm.synced_folder ".", "/home/ubuntu/go/src/github.com/puppetlabs/lumogon"
        m.vm.provider "virtualbox" do |vb|
          vb.memory = "2048"
        end
        m.vm.provision "shell", inline: provision_ubuntu(package)
      end
    end
  end

  machines['centos'].each do |distro, matrix|
    matrix['docker'].each do |version, package|
      config.vm.define "centos-#{distro}-#{version}" do |m|
        m.vm.box = matrix['box']
        m.vm.synced_folder ".", "/home/vagrant/go/src/github.com/puppetlabs/lumogon"
        m.vm.provider "virtualbox" do |vb|
          vb.memory = "2048"
        end
        m.vm.provision "shell", inline: provision_centos(distro, package)
      end
    end
  end
end

def provision_centos(distro, package)
  %{
  export GOPATH=/home/vagrant/go
  export PATH=$PATH:$GOPATH/bin

  puppet resource yumrepo docker \
         ensure=present \
         enabled=1 \
         gpgcheck=1 \
         gpgkey="https://yum.dockerproject.org/gpg" \
         baseurl="https://yum.dockerproject.org/repo/main/centos/#{distro}"

  puppet resource package docker-engine ensure=#{package}
  puppet resource package git ensure=present
  puppet resource service docker ensure=running
  if [ ! -d /usr/local/go ];
  then
    curl https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz -o /tmp/go1.8.3.linux-amd64.tar.gz
    tar -C /usr/local -xzf /tmp/go1.8.3.linux-amd64.tar.gz
  fi
  if [ ! -f /etc/profile.d/golang.sh ];
  then
    echo "export PATH=/usr/local/go/bin:$PATH" | tee /etc/profile.d/golang.sh && chmod +x /etc/profile.d/golang.sh
    source /etc/profile.d/golang.sh
  fi

  if [ ! -d /usr/local/rvm ];
  then
    yum install gcc-c++ patch readline readline-devel zlib zlib-devel
    yum install libyaml-devel libffi-devel openssl-devel make
    yum install bzip2 autoconf automake libtool bison iconv-devel sqlite-devel
    curl -sSL https://rvm.io/mpapis.asc | gpg --import -
    curl -L get.rvm.io | bash -s stable
    source /etc/profile.d/rvm.sh
    rvm reload
    rvm install 2.4.0
    rvm use 2.4.0 --default
    gem install bundler
  fi

  cd /home/vagrant/go/src/github.com/puppetlabs/lumogon && rm -rf vendor && make all

  docker version
  sudo docker run --rm \
       -v /var/run/docker.sock:/var/run/docker.sock \
       -e DOCKER_API_VERSION=1.19 \
       puppet/lumogon --disable-analytics version

  sudo docker run -d nginx
  sudo docker run --rm \
       -v /var/run/docker.sock:/var/run/docker.sock \
       -e DOCKER_API_VERSION=1.19 \
       puppet/lumogon --disable-analytics scan
}
end

def provision_ubuntu(package)
  %{
  export GOPATH=/home/ubuntu/go
  export PATH=$PATH:$GOPATH/bin
  sudo chown ubuntu:ubuntu -R $GOPATH
  sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
  echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" | tee /etc/apt/sources.list.d/docker.list
  sudo add-apt-repository ppa:longsleep/golang-backports
  sudo apt-get update
  sudo apt-get install -y --allow-downgrades golang-go docker-engine=#{package}
                                                       docker rmi puppet/lumogon
  cd /home/ubuntu/go/src/github.com/puppetlabs/lumogon && rm -rf vendor && make all

  docker version
  sudo docker run --rm \
       -v /var/run/docker.sock:/var/run/docker.sock \
       -e DOCKER_API_VERSION=1.24 \
       puppet/lumogon --disable-analytics version

  sudo docker run -d nginx
  sudo docker run --rm \
       -v /var/run/docker.sock:/var/run/docker.sock \
       -e DOCKER_API_VERSION=1.24 \
       puppet/lumogon --disable-analytics scan
}
end
