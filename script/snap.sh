#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

snapcraft upload --release=edge,beta,candidate,stable dist/red-cli_"$1"_linux_386.snap
snapcraft upload --release=edge,beta,candidate,stable dist/red-cli_"$1"_linux_amd64.snap
snapcraft upload --release=edge,beta,candidate,stable dist/red-cli_"$1"_linux_arm64.snap