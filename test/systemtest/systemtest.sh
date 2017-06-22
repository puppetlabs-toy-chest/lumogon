#!/usr/bin/env bash

###############################################################################
#   Constants
###############################################################################
DIND_VERSIONS=( 17.06.0 17.05.0 17.04.0 17.03.2 1.13.1 1.12.6 1.11.2 )
MISSIND_DIND_VERSIONS=( 1.10.3 1.9.1 1.8.3 1.7.1 1.6.2 1.5.0 )

IMAGE_NAME="puppet/lumogon"
TEST_IMAGE="nginx:stable-alpine"
OFFICIAL_DIND_IMAGE="docker"
MISSING_DIND_IMAGE="johnmccabe/dind"

TITLE="\033[104m" # white on light blue
PASS="\033[42m"
FAIL="\033[41m"
NC="\033[0m" # no colour

###############################################################################
#   Variables
###############################################################################
PASSED=()
FAILED=()
LUMODEBUG="${LUMODEBUG:-0}"
REDIRECT=""


###############################################################################
#   Functions
###############################################################################

function runLumogon()
{
    local dind_image=$1
    local dind_version=$2
    local dind_version_suffix=$3
    local current_dind_image="${dind_image}:${dind_version}${dind_version_suffix}"
    local scan_output=/tmp/dind-${dind_version}.lumoscan

    echo "============================================"
    echo -e "${TITLE}Docker version ${dind_version}${NC}"
    docker run --rm -d --privileged \
        --name dind-${dind_version} \
        ${dind_image}:${dind_version}${dind_version_suffix}

    docker run --rm \
        -e DOCKER_HOST=tcp://docker:2375 \
        -v ${PWD}/lumogon.tar:/lumogon.tar \
        --link dind-${dind_version}:docker \
        ${dind_image}:${dind_version}\
        docker load -i /lumogon.tar 

    docker run --rm -e \
        DOCKER_HOST=tcp://docker:2375 \
        -v ${PWD}/testimage.tar:/testimage.tar \
        --link dind-${dind_version}:docker \
        ${dind_image}:${dind_version} \
        docker load -i /testimage.tar

    docker run --rm \
        -e DOCKER_HOST=tcp://docker:2375 \
        --link dind-${dind_version}:docker \
        ${dind_image}:${dind_version}\
        docker run -d --name nginx ${TEST_IMAGE}

    rm -rf ${scan_output}

    docker run --rm \
        -e DOCKER_HOST=tcp://docker:2375 \
        --link dind-${dind_version}:docker \
        ${dind_image}:${dind_version} \
        docker run --rm -v /var/run/docker.sock:/var/run/docker.sock ${IMAGE_NAME} scan \
    > ${scan_output}

    docker kill dind-${v}
}

function saveImages()
{
    echo "Saving images to disk"
    echo "  - ${IMAGE_NAME}"
    docker save ${IMAGE_NAME} > ${PWD}/lumogon.tar

    echo "  - ${TEST_IMAGE}"
    docker pull ${TEST_IMAGE}
    docker save ${TEST_IMAGE} > ${PWD}/testimage.tar

}

function checkNumContainers()
{
    local dind_version=$1
    local scan_output=/tmp/dind-${dind_version}.lumoscan

    if [ ! -f ${scan_output} ]; then
        echo -e "  - ${FAIL}no scan file found${NC}"
        return 1
    fi

    local numContainers=`jq '.containers | length' ${scan_output}`
    if [ "$numContainers" != "1" ]; then
        echo -e "  - ${FAIL}unexpected number of containers returned [${numContainers}]${NC}"
        return 1
    fi

    echo -e "  - ${PASS}expected number of containers returned 1${NC}"
}


###############################################################################
#   Main
###############################################################################

saveImages

echo "Running ${IMAGE_NAME} against official Docker DIND images"
for v in "${DIND_VERSIONS[@]}"
do
    runLumogon ${OFFICIAL_DIND_IMAGE} ${v} -dind
    checkNumContainers ${v}
done

echo "Running ${IMAGE_NAME} against Missing DIND images"
for v in "${MISSIND_DIND_VERSIONS[@]}"
do
    runLumogon ${MISSING_DIND_IMAGE} ${v}
done

echo "Finished"
