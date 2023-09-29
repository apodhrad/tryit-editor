#!/bin/bash
IMAGE="quay.io/apodhrad/tryit-editor:latest"
podman run -it --rm -p 8080:8080 -v ./custom:/var/tryit-editor:Z "${IMAGE}" start -c /var/tryit-editor/services.yaml
