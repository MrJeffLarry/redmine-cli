#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

#git push --delete origin "v$version"
#git tag --delete "v$version"

git tag -a v$1 -m "Release v$1"
git push origin v$1