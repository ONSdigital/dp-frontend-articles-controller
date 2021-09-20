#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-frontend-articles-controller
  make audit
popd