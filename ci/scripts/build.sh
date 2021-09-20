#!/bin/bash -eux

pushd dp-frontend-articles-controller
  make build
  cp build/dp-frontend-articles-controller Dockerfile.concourse ../build
popd
