#!/bin/bash

# Login to Docker
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

# Install goreleaser
go get -u github.com/goreleaser/goreleaser

# Create the release
goreleaser release
