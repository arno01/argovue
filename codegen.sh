#!/bin/bash 

mkdir -p /tmp/go/src
export GOPATH=/tmp/go
go get k8s.io/code-generator
go get k8s.io/apimachinery

ln -s $(pwd)/src $GOPATH/src/argovue

$GOPATH/src/k8s.io/code-generator/generate-groups.sh all argovue/client argovue/apis "kubevue.io:v1"

rm -rf /tmp/go
