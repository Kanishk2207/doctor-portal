#!/bin/bash
source config.sh

echo "Usage: ./build <dockerfile>"

if docker build -f $DOCKER_FILE -t $IMAGE_NAME:$IMAGE_TAG .; then
    echo "Image built successfully"
    exit 0
else
    echo "Image building failed, exit status: $?"
    exit $?
fi


