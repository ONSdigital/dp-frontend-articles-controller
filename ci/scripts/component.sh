#!/bin/bash -eux

pushd dp-frontend-articles-controller
  make test-component
popd
