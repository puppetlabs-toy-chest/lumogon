# Lumogon Client - lumogon

**NOTE** All naming is placeholder and tbd.

## Overview
TODO

## Building

### Prerequisites

This assumes you have Go 1.8 installed on your system with your $GOPATH configured.

When working from a fork you should avoid using `go get` to grab the fork and instead either clone the repo into the `$GOPATH/src/github.com/puppetlabs/` directory or add your fork as a new remote and switch to that.

```
go get github.com/puppetlabs/transparent-containers
cd $GOPATH/github.com/puppetlabs/transparent-containers
git remote rename origin upstream
git remote add origin git@github.com:<your fork>/transparent-containers.git
```

### Building a local Lumogon container
You can build the Lumogon client and add it to a scratch container image by running:
```
make bootstrap
make all
```
This will create an image called `local/lumogon` which you can run as follows (passing additional commands as required):
```
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock local/lumogon help
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock local/lumogon version
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock local/lumogon scan <target container ID or name>
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock local/lumogon capability list
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock local/lumogon capability describe <capability name>
```
General flags, there is currenly one main general flag:
| Long Flag     | Short Flag    | Effect                                                         |
| ------------- |:-------------:|:--------------------------------------------------------------:|
| `--debug`     | `-d`          | Enables debug logging to stdout, this is suppressed by default |

### Building a Lumogon OSX client (optional)
The provided `Makefile` compiles a binary for use within a `SCRATCH` container, if you want to build a client to run on OSX you can override the default OS as follows:
```
GOOS=darwin make test build
```
This will build an OSX compatible binary in `bin/lumogon`, this is super useful for local testing and gives you the same functionality as running in the container.

## Running
Currently running `lumogon container <target container ID or name>` (via Docker container or OSX binary) will spin up a harvesting container which gathers capability data and dumps it to stdout, or if the `--consumer` flag is set, to the supplied consumer endpoint.

### Examples
**Harvesting all running containers**
```
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock local/lumogon scan
```

**Keeping a harvester container running after gathering data**
```
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock local/lumogon scan nginx --keep-harvesters
```

**Submitting harvested data to a custom endpoint**
```
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock local/lumogon report redis --endpoint http://domain/my/endpoint/
```

**Showing full output, both stdout and stderr**
```
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock local/lumogon scan 3361d2d49d5a
[lumogon] 2017/04/04 09:35:11.196250 Initialising capability: Host Capability
[lumogon] 2017/04/04 09:35:11.196272 Registering capability: Host Capability
[lumogon] 2017/04/04 09:35:11.196276 Initialising capability: Label Capability
[lumogon] 2017/04/04 09:35:11.196280 Registering capability: Label Capability
[lumogon] 2017/04/04 09:35:11.197079 Creating container runtime client: docker
[lumogon] 2017/04/04 09:35:11.199443 Image already downloaded: local/lumogon
[lumogon] 2017/04/04 09:35:11.199498 Harvesting capabilities from [1] containers
[lumogon] 2017/04/04 09:35:11.200894 Creating harvester container: lumogon_wonderful_curie8
[lumogon] 2017/04/04 09:35:11.200917 Attaching harvester to container ID: 3361d2d49d5a4e35b4e5eb7aee723a5c47c3f59c3d1f63620a7b593ad3b6151c, Name: /autopilotpatternhelloworld_nginx_1
[lumogon] 2017/04/04 09:35:12.032244 Started harvester container, ID: 16d8bedde570cb1a6408381307906d7a17028b29c481997ea42a32a2f8165452
[lumogon] 2017/04/04 09:35:12.033256 Harvesting complete
[lumogon] 2017/04/04 09:35:12.764089 Removed harvester container: 16d8bedde570cb1a6408381307906d7a17028b29c481997ea42a32a2f8165452
{"$schema":"http://puppet.com/lumogon/core/draft-01/schema#1","generated":"2017-04-04 09:35:11.192957005 +0000 UTC","owner":"default","group":["default"],"build_information":{"version":"development","type":"","git_sha":"b44e0e0eeb8f80e9815adda636fe687b8404b479","built":"2017-04-04 09:24:44 UTC"},"reportid":"3f166706-7e1d-49db-9d58-6b74e745fa65","containers":{"3361d2d49d5a4e35b4e5eb7aee723a5c47c3f59c3d1f63620a7b593ad3b6151c":{"$schema":"http://puppet.com/lumogon/containerrecord/draft-01/schema#1","generated":"2017-04-04 09:35:11.966796229 +0000 UTC","container_report_id":"2931dd4c-af3b-4118-b026-95925e0c155c","container_id":"3361d2d49d5a4e35b4e5eb7aee723a5c47c3f59c3d1f63620a7b593ad3b6151c","container_name":"/autopilotpatternhelloworld_nginx_1","capabilities":{"host":{"$schema":"http://puppet.com/lumogon/capability/host/draft-01/schema#1","title":"Host Capability","harvestid":"732641e7-799d-48ae-9ecf-b8f463b2c1a8","payload":{"bootTime":1491121815,"hostid":"25239174-ded5-4e3c-9835-1c14b5cf3dae","hostname":"16d8bedde570","kernelVersion":"4.9.13-moby","os":"linux","platform":"alpine","platformFamily":"alpine","platformVersion":"3.4.4","procs":62,"uptime":176696,"virtualizationRole":"guest","virtualizationSystem":"docker"}},"label":{"$schema":"http://puppet.com/lumogon/capability/label/draft-01/schema#1","title":"Label Capability","harvestid":"2f221365-24da-43d1-be30-a3233c03bb4e","payload":{"com.docker.compose.config-hash":"ade07648dd8a8e46566e280385c7781f78d343d019a10fbb7200ee4cc87d35d5","com.docker.compose.container-number":"1","com.docker.compose.oneoff":"False","com.docker.compose.project":"autopilotpatternhelloworld","com.docker.compose.service":"nginx","com.docker.compose.version":"1.11.1"}}}}}}
>>>>>>> db61eca10121b915bc0d2335768509ca2e742114
```

**Showing only stdout (useful to pipe to other commands, jq etc.)**
```
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock local/lumogon scan autopilotpatternhelloworld_nginx_1 2>/dev/null
{"$schema":"http://puppet.com/lumogon/core/draft-01/schema#1","generated":"2017-04-04 09:36:16.544659452 +0000 UTC","owner":"default","group":["default"],"build_information":{"version":"development","type":"","git_sha":"b44e0e0eeb8f80e9815adda636fe687b8404b479","built":"2017-04-04 09:24:44 UTC"},"reportid":"d8c390a0-58be-4de9-8c4d-31746ca5618d","containers":{"3361d2d49d5a4e35b4e5eb7aee723a5c47c3f59c3d1f63620a7b593ad3b6151c":{"$schema":"http://puppet.com/lumogon/containerrecord/draft-01/schema#1","generated":"2017-04-04 09:36:17.324595111 +0000 UTC","container_report_id":"9823ce2e-0169-44fe-86f7-fe1d1bd9d15a","container_id":"3361d2d49d5a4e35b4e5eb7aee723a5c47c3f59c3d1f63620a7b593ad3b6151c","container_name":"/autopilotpatternhelloworld_nginx_1","capabilities":{"host":{"$schema":"http://puppet.com/lumogon/capability/host/draft-01/schema#1","title":"Host Capability","harvestid":"d81dce4c-b062-4695-a6fe-1e195b729da4","payload":{"bootTime":1491121815,"hostid":"25239174-ded5-4e3c-9835-1c14b5cf3dae","hostname":"84a12473fd44","kernelVersion":"4.9.13-moby","os":"linux","platform":"alpine","platformFamily":"alpine","platformVersion":"3.4.4","procs":62,"uptime":176762,"virtualizationRole":"guest","virtualizationSystem":"docker"}},"label":{"$schema":"http://puppet.com/lumogon/capability/label/draft-01/schema#1","title":"Label Capability","harvestid":"c73d9647-1eb0-4b92-9425-39c0af229859","payload":{"com.docker.compose.config-hash":"ade07648dd8a8e46566e280385c7781f78d343d019a10fbb7200ee4cc87d35d5","com.docker.compose.container-number":"1","com.docker.compose.oneoff":"False","com.docker.compose.project":"autopilotpatternhelloworld","com.docker.compose.service":"nginx","com.docker.compose.version":"1.11.1"}}}}}}
```
