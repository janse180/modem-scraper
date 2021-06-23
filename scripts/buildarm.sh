#!/bin/bash

BUILD_VERSION="$(date +'%Y%m%d%H%M%S')-$(git log --format=%h -1)"

echo "Building arm binary with version [${BUILD_VERSION}]..."

BINARY_NAME="modem-scraper"
FLAGS="-X main.BuildVersion=${BUILD_VERSION}"
GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="${FLAGS}" -o ${BINARY_NAME}

ARCHIVE_NAME="${BINARY_NAME}_${BUILD_VERSION}_linux_arm64.zip"
zip ${ARCHIVE_NAME} ${BINARY_NAME}

if [ "${TRAVIS_BRANCH}" = "master" ] && [ "${TRAVIS_PULL_REQUEST_BRANCH}" = "" ] ; then

  echo "Building image..."
  IMAGE="janse180/modem-scraper"
  docker build . -t ${IMAGE}:${BUILD_VERSION}-arm64

  echo "Pushing image to registry..."
  echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin ${DOCKER_URL}
  docker push ${IMAGE}:${BUILD_VERSION}-arm64
fi