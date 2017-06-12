bootstrap:
	sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	echo "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(shell lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list
	sudo apt-get update
	sudo apt-get install docker-ce

dependencies:
	bundle install
	docker pull docker:dind
	docker pull docker:latest
	docker pull puppet/lumogon:latest

test: dependencies
	bundle exec rake test

deploy:
	true

all: dependencies test

.PHONY: dependencies bootstrap test deploy all
