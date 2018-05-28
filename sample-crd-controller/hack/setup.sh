#!/usr/bin/env bash

set -xe

GOPATH=$(go env GOPATH)
PACKAGE_NAME=k8s-practice/sample-crd-controller
REPO_ROOT="$GOPATH/src/github.com/shudipta/$PACKAGE_NAME"

pushd $REPO_ROOT

go build -o hack/docker/sample-crd-controller main.go
chmod +x hack/docker/sample-crd-controller

docker build -t shudipta/samp-crd-ctl:latest hack/docker
docker save shudipta/samp-crd-ctl:latest | pv | (eval $(minikube docker-env) && docker load)
#docker push shudipta/samp-crd-ctl:latest

rm -rf hack/docker/sample-crd-controller
popd