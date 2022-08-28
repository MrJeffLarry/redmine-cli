#!/bin/bash

export EMAIL=jeff@hagerman.io

rm ./apt/*.deb

cp ./dist/*.deb ./apt/

cd apt

# Packages & Packages.gz
dpkg-scanpackages --multiversion . > Packages
gzip -k -f Packages

# Release, Release.gpg & InRelease
apt-ftparchive release . > Release
gpg --default-key "$EMAIL" -abs -o - Release > Release.gpg
gpg --default-key "$EMAIL" --clearsign -o - Release > InRelease