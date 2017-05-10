# What about Images?

The current tooling focuses on scanning containers, but you may want to
scan an image instead. Luckily an image is just a container, albeit
probably one that's not running.

So to scan an image we can write a simle script which boots an image and
then runs `lumogon` agains the image instance. For instance:

```
./scan-image.sh ocaml/opam
```

This will pull the image (in this case `ocaml/opam`) and output to
stdout the Lumogon JSON output. Let's use JQ to find out the operating
system details of a third party image from Docker Hub:

```
$ ./scan-image.sh elasticsearch | jq ".containers[].capabilities.host.payload"
{
  "hostname": "55b3031db6b4",
  "kernelversion": "4.9.21-moby",
  "os": "linux",
  "platform": "debian",
  "platformfamily": "debian",
  "platformversion": "8.7",
  "procs": "61",
  "uptime": "263071",
  "virtualizationrole": "guest",
  "virtualizationsystem": "docker"
}
```

## Augmenting images with Lumogon

Another trick is to embed the data generated from Lumogon inside an
image. Here we:

* Pull the named image if it's not available locally
* Run an instance of the image
* Use Lumogon to conduct a scan
* Copy the scan output the the running container
* Save a new version of the image with the scan output on a new layer

```
./augment-image.sh <your/image>
```

This has some interesting outcomes. For instance you can now access the
data from Lumogon incredibly quickly via `docker run` or via `docker exec`. eg:

```
docker run --rm <your/image> cat /lumogon.json
```

Not only is this faster that scanning each time but if you're using
immutable containers via the `--read-only` flag then you have some
guarentees the data hasn't been tampered with.


You'll find the code for these examples in the [Lumogon
repository](https://github.com/puppetlabs/lumogon/tree/master/examples).
Both of these examples demonstrate how easy it is to built on top of
Lumogon and point the direction for the kinds of features we might add
into the tool at a later stage.

If you're interested in these examples or have any questions about Lumogon
then head over to our [Slack channel](https://puppetcommunity.slack.com/messages/G58F97FC5).
