#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-s3
  make build
popd

