#!/usr/bin/env bash

# Builds a docker image for the application, and then pushes it to a remote server
# for deployment.
#
# To set the IP address of the remote server to push the docker image to, create a
# configuration file named build-and-push.config in the same directory, containing
# the line REMOTE_IP=<your server IP>. Check out the build-and-push.config.example
# file as an example.
#
# Usage:
# ./build-and-push.sh <docker image tag>
#
# Example usage:
# ./build-and-push.sh betonz:1

source build-and-push.config

TAG=$1
IMAGE_TARBALL=$(echo $TAG | tr : -).tar.gz
COPY_DESTINATION=root@$REMOTE_IP:/images/.

echo Building docker image $TAG
docker build . -t $1 || { echo Build failed. Aborting! ; exit 1; }

echo Saving image to $IMAGE_TARBALL
docker save $1 | gzip > $IMAGE_TARBALL

echo Copying to $COPY_DESTINATION
scp $IMAGE_TARBALL $COPY_DESTINATION
rm $IMAGE_TARBALL

echo Loading docker image in remote
ssh root@$REMOTE_IP -t "docker image rm $TAG ; docker load < /images/$IMAGE_TARBALL"

echo Done!
