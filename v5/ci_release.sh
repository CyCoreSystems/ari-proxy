#!/bin/bash

# Only run if we are inside Travis-CI
if [ ! -e $CI ]; then
   echo "Logging in to Docker..."
   echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

   # Create the release
   echo "Creating release..."
   goreleaser release
fi

