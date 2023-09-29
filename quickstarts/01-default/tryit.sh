#!/bin/bash
IMAGE="quay.io/apodhrad/tryit-editor:latest"
podman run -it --rm -p 8080:8080 "${IMAGE}"
