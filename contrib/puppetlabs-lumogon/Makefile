bootstrap:
	apt-get install apt-transport-https ca-certificates curl software-properties-common
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
	echo "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt.sources.d/docker.list
	apt-get update
	apt-get install docker-ce

dependencies:
	bundle install

test: dependencies
	bundle exec rake test

.PHONY: dependencies bootstrap test
