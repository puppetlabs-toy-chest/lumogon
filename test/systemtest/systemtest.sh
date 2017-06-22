#!/usr/bin/env bash

DIND_VERSIONS=( 17.06.0 17.05.0 17.04.0 17.03.2 1.13.1 1.12.6 1.11.2 )
MISSIND_DIND_VERSIONS=( 1.10.3 1.9.1 1.8.3 1.7.1 1.6.2 1.5.0 )

IMAGE_NAME="puppet/lumogon"
TEST_IMAGE="nginx:stable-alpine"
OFFICIAL_DIND_IMAGE="docker"
MISSING_DIND_IMAGE="johnmccabe/dind"

TITLE="\033[104m" # white on light blue
NC="\033[0m" # no colour

echo "Saving local image ${IMAGE_NAME} to disk..."
docker save ${IMAGE_NAME} > ${PWD}/lumogon.tar

echo "Saving test image ${TEST_IMAGE} to disk..."
docker pull ${TEST_IMAGE}
docker save ${TEST_IMAGE} > ${PWD}/testimage.tar

echo "Running ${IMAGE_NAME} against official Docker DIND images"
for v in "${DIND_VERSIONS[@]}"
do
    echo -e "  - ${TITLE}Docker version ${v}${NC}"
    docker run --rm -d --privileged --name dind-${v} docker:${v}-dind
    docker run --rm -e DOCKER_HOST=tcp://docker:2375 -v ${PWD}/lumogon.tar:/lumogon.tar --link dind-${v}:docker ${OFFICIAL_DIND_IMAGE}:${v} docker load -i /lumogon.tar
    docker run --rm -e DOCKER_HOST=tcp://docker:2375 -v ${PWD}/testimage.tar:/testimage.tar --link dind-${v}:docker ${OFFICIAL_DIND_IMAGE}:${v} docker load -i /testimage.tar
    docker run --rm -e DOCKER_HOST=tcp://docker:2375 --link dind-${v}:docker ${OFFICIAL_DIND_IMAGE}:${v} docker run -d --name nginx ${TEST_IMAGE}
    docker run --rm -e DOCKER_HOST=tcp://docker:2375 --link dind-${v}:docker ${OFFICIAL_DIND_IMAGE}:${v} docker run --rm -v /var/run/docker.sock:/var/run/docker.sock ${IMAGE_NAME} scan
    docker kill dind-${v}
done

echo "Running ${IMAGE_NAME} against Missing DIND images"
for v in "${MISSIND_DIND_VERSIONS[@]}"
do
    echo -e "  - ${TITLE}Docker version ${v}${NC}"
    docker run --rm -d --privileged --name dind-${v} ${MISSING_DIND_IMAGE}:${v}
    docker run --rm -e DOCKER_HOST=tcp://docker:2375 -v ${PWD}/lumogon.tar:/lumogon.tar --link dind-${v}:docker ${MISSING_DIND_IMAGE}:${v} docker load -i /lumogon.tar
    docker run --rm -e DOCKER_HOST=tcp://docker:2375 -v ${PWD}/testimage.tar:/testimage.tar --link dind-${v}:docker ${OFFICIAL_DIND_IMAGE}:${v} docker load -i /testimage.tar
    docker run --rm -e DOCKER_HOST=tcp://docker:2375 --link dind-${v}:docker ${MISSING_DIND_IMAGE}:${v} docker run -d --name nginx ${TEST_IMAGE}
    docker run --rm -e DOCKER_HOST=tcp://docker:2375 --link dind-${v}:docker ${MISSING_DIND_IMAGE}:${v} docker run --rm -v /var/run/docker.sock:/var/run/docker.sock ${IMAGE_NAME} scan
    docker kill dind-${v}
done

echo "Finished"