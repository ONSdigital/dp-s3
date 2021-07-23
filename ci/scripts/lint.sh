#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-s3
  make lint
popd
