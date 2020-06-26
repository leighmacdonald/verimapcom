#!/bin/bash
IMAGE_ID=$(docker images | head -2 | tail -1 | awk '{print $3}')
docker tag "${IMAGE_ID}" docker.pkg.github.com/leighmacdonald/verimapcom/verimapcom:latest
docker push docker.pkg.github.com/leighmacdonald/verimapcom/verimapcom:latest