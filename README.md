# Lumogon

[![Build
Status](https://travis-ci.com/puppetlabs/lumogon.svg?token=RqtxRv25TsPVz69Qso5L&branch=master)](https://travis-ci.com/puppetlabs/lumogon)

## Introduction

The Lumogon tool provides a way to inspect the
[Docker](https://www.docker.com/) containers you have running on your computer.
`lumogon` can produce a report of what's running, along with the various things
happening inside those containers.

Taking this a step further, Lumogon also works with a website hosted on the internet
where you can send reports to view online, and to share with others.

## Trying out Lumogon

### Getting access to the client software

You'll need [Docker](https://www.docker.com/) installed and running locally.
This should already be true if you have Docker containers you want to analyze
with Lumogon.

You can get the Lumogon client via one of two methods:

**Downloading the client from Docker Hub:**

**NOTE: we do not yet have an image on Docker Hub. Until such time you will want to build the client from source. See below.**

This is the simplest way to get the Lumogon client. At a terminal, do the following:

``` shell
docker pull puppetlabs/lumogon
```

**Building the client from source:**

For this you'll need a few more things:

 - Install [Go](https://golang.org/dl/), version 1.8 or later
 - Set your `$GOPATH` variable to the path where you want to keep your Go sources -- for example, `${HOME}/go`.
 - Download the Lumogon project
 - Build the cli Docker image

Now, at a terminal type:

``` shell
export GOPATH="${HOME}/go"
mkdir -p ${GOPATH}/src/github.com/puppetlabs
cd ${GOPATH}/src/github.com/puppetlabs
git clone https://github.com/puppetlabs/lumogon
cd $GOPATH/src/github.com/puppetlabs/lumogon
make all
```

### Running the client against your local docker containers

Now that you have the Lumogon client installed, let's run it to find out which
containers you have running and what we can learn about them. The output of the
client will be a [JSON](https://en.wikipedia.org/wiki/JSON) listing of all the
containers found and what Lumogon could learn about them.

Since JSON can be hard to read as one big unstructured message, you might want
to use a "pretty-printer" like [jq](https://stedolan.github.io/jq/) to print the
data in a more readable format.

``` shell
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon scan | jq .
```

After a few seconds you should see your JSON data:

``` json
{
  "$schema": "http://puppet.com/lumogon/core/draft-01/schema#1",
  "generated": "2017-04-18 22:30:21.652393839 +0000 UTC",
  "owner": "default",
  "group": [
    "default"
  ],
  "client_version": {
    "BuildVersion": "development",
    "BuildTime": "2017-04-18 08:50:26 UTC",
    "BuildSHA": "7983b886af060dcaba171ebae393e2b31ff57063"
  },
  "reportid": "a862f238-e817-42fc-a972-0563c2cc3992",
  "containers": {
    "fd9a36bbc94cc0503fa87134fb48eb6966f6e70cd663764fc6b201a91dd90d6a": {
      "$schema": "http://puppet.com/lumogon/containerreport/draft-01/schema#1",
      "generated": "2017-04-18 22:30:20.584745556 +0000 UTC",
      "container_report_id": "2a92817d-186f-47e6-9d30-e0df2cdfbf89",
      "container_id": "fd9a36bbc94cc0503fa87134fb48eb6966f6e70cd663764fc6b201a91dd90d6a",
      "container_name": "/dynamodb",
      "capabilities": {
        "dpkg": {
          "$schema": "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
          "title": "Dpkg Capability",
          "type": "dockerapi",
          "harvestid": "09401dde-1b8f-47d7-83b5-ab525f5539eb",
          "payload": {
            "packages": [
              "acl,2.2.52-2",
              "adduser,3.113+nmu3",
              "apt,1.0.9.8.4",
              "base-files,8+deb8u7",
              "base-passwd,3.5.37",
              "bash,4.3-11+deb8u1",
              "bsdutils,1:2.25.2-6",
              "bzip2,1.0.6-7+b3",
              "...skip-a-bit...,0.0.0",
              "zlib1g,1:1.2.8.dfsg-2+b1"
            ]
          }
        },
        "host": {
          "$schema": "http://puppet.com/lumogon/capability/host/draft-01/schema#1",
          "title": "Host Capability",
          "type": "attached",
          "harvestid": "226fba7a-a913-45f9-8f46-d57d719d30a5",
          "payload": {
            "BootTime": 1492180953,
            "HostID": "a2cbff83-c8ad-e238-b8f4-59665c5c1e34",
            "Hostname": "e740aad631fb",
            "KernelVersion": "4.9.13-moby",
            "OS": "linux",
            "Platform": "debian",
            "PlatformFamily": "debian",
            "PlatformVersion": "8.7",
            "Procs": 61,
            "Uptime": 373667,
            "VirtualizationRole": "guest",
            "VirtualizationSystem": "docker"
          }
        },
        "label": {
          "$schema": "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
          "title": "Label Capability",
          "type": "dockerapi",
          "harvestid": "dc50f42f-f459-4100-8c09-db3a1ef88335",
          "payload": {
            "default": "No labels found"
          }
        }
      }
    }
  }
}
```

Since this output is valid JSON, you can slice and dice it with `jq`, or pass it
along to any other tool you use that can accept JSON input.

### Sending container reports to the reporting website

If you provide a `--endpoint` argument with a URL to the lumogon program,
it will send your JSON data to our reporting application where you can share the
reports via a URL.

Try this:

``` shell
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon report --endpoint https://consumer.api.lumogon.com/api/v1/
```

If all goes as planned, you should see a line like the following with a URL you can visit to view your results:

```
http://reporter.api.lumogon.com/aHrhCcXT2sJBrrewEFCGaWEYbJnqV0vQWMwzO3Dzhbc=
```

### More ways to use the Lumogon client

Let's get the usage for the Lumogon client:

``` shell
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon --help
Creates and attaches a container to the target, which harvests data from the target and sends to the Lumogon service

Usage:
  lumogon [command]

Available Commands:
  help        Help about any command
  report      Scan one or more containers and send the collected information to the Lumogon service
  scan        Scan one or more containers and print the collected information
  version     Lumogon client version

Flags:
      --config string     config file (default is $HOME/.cli.yaml)
  -d, --debug             Print debug logging
  -k, --keep-harvesters   Keeps harvester containers instead of automatically deleting
```

Feel free to explore those command-line options.  Of note:

 - The "--keep-harvesters" flag will preserve containers that are created on the fly to explore your other containers. You can use `docker logs <containerid>" to see more of what they found while running.
 - You can specify `scan` to collect data on all your running containers, or you can specify a specific container by passing `scan <containerid>`.
 - `-d` will generate more verbose debugging output so you can see Lumogon exploring your containers.

### Giving us feedback

We'd love to hear from you. Feel free to contact us via email at (**TODO: set up the contact email address**). We also
run a [Slack](https://slack.com/) channel at (**TODO: Set up the Slack channel**). And, if you're a developer
and want to know more about how all this works, check out the next section.
