# Lumogon

[![Build
Status](https://travis-ci.com/puppetlabs/lumogon.svg?token=RqtxRv25TsPVz69Qso5L&branch=master)](https://travis-ci.com/puppetlabs/lumogon)

Lumogon provides a way to inspect, analyze and report on your running
[Docker](https://www.docker.com/) containers.

## Trying out Lumogon

### Downloading Lumogon

You'll need [Docker](https://www.docker.com/) installed and running locally.
This should already be true if you have Docker containers you want to analyze
with Lumogon.

``` shell
docker pull puppet/lumogon
```

### Running Lumogon

Now that you have Lumogon installed, let's run it to find out which
containers you have running and what we can learn about them. The output of the
client will be a [JSON](https://en.wikipedia.org/wiki/JSON) listing of all the
containers found and what Lumogon could learn about them.

Since JSON can be hard to read as one big unstructured message, you might want
to use a "pretty-printer" like [jq](https://stedolan.github.io/jq/) to print the
data in a more readable format.

``` shell
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon scan | jq
```

After a few seconds you should see your JSON data:

``` json
{
  "$schema": "http://puppet.com/lumogon/core/draft-01/schema#1",
  "generated": "2017-05-09 07:59:24.287008012 +0000 UTC",
  "owner": "default",
  "group": [
    "default"
  ],
  "client_version": {
    "BuildVersion": "development",
    "BuildTime": "2017-05-09 06:56:22 UTC",
    "BuildSHA": "9e8f684432ff12b04b5b5d594caa0ebcce86b844"
  },
  "reportid": "c73a79dc-8612-4af8-8bd8-22e32ea11e38",
  "containers": {
    "5982d3f16bbbf9530ae09915b22a0d189044e3b953e5e417e2783b90de579034": {
      "$schema": "http://puppet.com/lumogon/containerreport/draft-01/schema#1",
      "generated": "2017-05-09 07:59:03.513739277 +0000 UTC",
      "container_report_id": "8d17e541-11b3-4f25-b145-4ad9d3045995",
      "container_id": "5982d3f16bbbf9530ae09915b22a0d189044e3b953e5e417e2783b90de579034",
      "container_name": "/fixtures_alpine_1",
      "capabilities": {
        "apk": {
          "$schema": "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
          "title": "Packages (APK)",
          "type": "dockerapi",
          "harvestid": "3a5bf0d4-36d8-440b-af81-615b5493fe98",
          "payload": {
            "alpine-baselayout": "3.0.3-r0",
            "alpine-keys": "1.1-r0",
            "apk-tools": "2.6.7-r0",
            "busybox": "1.24.2-r9",
            "libc-utils": "0.7-r0",
            "libcrypto1.0": "1.0.2h-r1",
            "libssl1.0": "1.0.2h-r1",
            "musl": "1.1.14-r10",
            "musl-utils": "1.1.14-r10",
            "scanelf": "1.1.6-r0",
            "zlib": "1.2.8-r2"
          }
        },
        "dpkg": {
          "$schema": "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
          "title": "Packages (DPKG)",
          "type": "dockerapi",
          "harvestid": "bdee3efe-70cf-4684-9eb6-cfbfeeb96b9c"
        },
        "host": {
          "$schema": "http://puppet.com/lumogon/capability/host/draft-01/schema#1",
          "title": "Host Information",
          "type": "attached",
          "harvestid": "53d1961c-e8e9-4b52-8620-6bac37a69664",
          "payload": {
            "hostname": "365cfca386ec",
            "kernelversion": "4.9.21-moby",
            "os": "linux",
            "platform": "alpine",
            "platformfamily": "alpine",
            "platformversion": "3.4.0",
            "procs": "61",
            "uptime": "248396",
            "virtualizationrole": "guest",
            "virtualizationsystem": "docker"
          }
        },
        "label": {
          "$schema": "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
          "title": "Labels",
          "type": "dockerapi",
          "harvestid": "50a3f846-3580-4190-9b00-27c3011f1516",
          "payload": {
            "com.docker.compose.config-hash": "70e9897635135adc7e9bd0af535fef48ae8e26c8e0debbf8f40e0d67938a9884",
            "com.docker.compose.container-number": "1",
            "com.docker.compose.oneoff": "False",
            "com.docker.compose.project": "fixtures",
            "com.docker.compose.service": "alpine",
            "com.docker.compose.version": "1.11.2"
          }
        },
        "rpm": {
          "$schema": "http://puppet.com/lumogon/capability/label/draft-01/schema#1",
          "title": "Packages (RPM)",
          "type": "dockerapi",
          "harvestid": "1c2976f2-802c-4138-a6b7-e7a814340fea"
        }
      }
    }
  }
}
```

Since this output is valid JSON, you can slice and dice it with `jq`, or pass it
along to any other tool you use that can accept JSON input.

### Sending reports to the Lumogon service

As well as being able to access raw data we're building a web service
which can provide more visual reports.


``` shell
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon report
```

This collects the same data as `scan` but sends it over HTTPS to our
service. Instead of returning the raw data you should get a URL to view
the report, like this one:

```
http://reporter.app.lumogon.com/aHrhCcXT2sJBrrewEFCGaWEYbJnqV0vQWMwzO3Dzhbc=
```

### More ways to use Lumogon

Let's get the usage for the Lumogon client:

``` shell
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon --help
Lumogon is a tool for inspecting, reporting on, and analyzing your container applications.

Usage:
  lumogon [command]

Available Commands:
  help        Help about any command
  report      Scan one or more containers and send the collected information to the Lumogon service
  scan        Scan one or more containers and print the collected information
  version     Show the Lumagon version information

Flags:
  -d, --debug             Print debug logging
  -k, --keep-harvesters   Keeps harvester containers instead of automatically deleting
```

Feel free to explore those command-line options. Of note:

 - The `--keep-harvesters` flag will preserve containers that are created on the fly to explore your other containers. You can use `docker logs <containerid>` to see more of what they found while running.
 - You can specify `scan` to collect data on all your running containers, or you can specify a specific container by passing `scan <containerid>`.
 - `--debug` will generate more verbose debugging output so you can see Lumogon exploring your containers.


## Building the client from source

If you're making changes to Lumogon, or just interested in seeing how it works under the hood, you might want to try building from source. For this you'll need a few more things:

 - Install [Go](https://golang.org/dl/), version 1.8 or later
 - Install `make`
 - Set your `$GOPATH` variable to the path where you want to keep your Go sources -- for example, `${HOME}/go`.
 - Download the Lumogon source code
 - Build the Docker image

The terminal commands to do this are:

```shell
export GOPATH="${HOME}/go"
mkdir -p ${GOPATH}/src/github.com/puppetlabs
cd ${GOPATH}/src/github.com/puppetlabs
git clone https://github.com/puppetlabs/lumogon
cd $GOPATH/src/github.com/puppetlabs/lumogon
make all
```

Note that this build process isn't widely tested away from macOS just yet but will eventually work everywhere.


## Giving us feedback

We'd love to hear from you. We have a [Slack channel](https://puppetcommunity.slack.com/messages/G58F97FC5) for talking about Lumogon and please do open issues against [the repository](https://github.com/puppetlabs/lumogon/issues).
