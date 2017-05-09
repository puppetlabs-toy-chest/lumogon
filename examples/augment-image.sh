#!/bin/bash -e

# This script augments a specified Docker image with a static inventory of information
# It does this by booting a container based on the image, generating the inventory and
# then saving a new version of the image with an additional layer containing the inventory
# at /lumogon.json
# This allows for accessing the inventory in any instance of that image, for instance:
# docker run --rm <image> cat /lumogon.json

IMAGE=${1?"Usage: $0 IMAGE[:TAG]"}

CONTAINER_ID=$(docker run --rm -di --entrypoint "" "$IMAGE" /bin/sh -c "while true; do echo hello world; sleep 1; done")
docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock local/tc container "$CONTAINER_ID" > output.json
docker cp ./output.json "$CONTAINER_ID":/lumogon.json
docker commit "$CONTAINER_ID" "$IMAGE"
docker kill "$CONTAINER_ID"
rm ./output.json
