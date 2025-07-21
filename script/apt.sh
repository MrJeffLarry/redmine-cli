#!/bin/bash

set -e

cp ./dist/*.deb ./apt/

cd ./apt

# Packages & Packages.gz
dpkg-scanpackages --multiversion . > Packages

gzip -k -f Packages

# Release, Release.gpg & InRelease
apt-ftparchive release . > Release
gpg --batch --pinentry-mode loopback --default-key "$GPG_FINGERPRINT" -abs -o - Release > Release.gpg
gpg --batch --pinentry-mode loopback --default-key "$GPG_FINGERPRINT" --clearsign -o - Release > InRelease

# git commit and push
git config --local user.name "GitHub Actions"
git config --local user.email "<>"
git pull
git add .
git commit -m "Update apt files"
git push origin main