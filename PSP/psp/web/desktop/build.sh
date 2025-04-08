#!/usr/bin/env bash
FE=$(pwd)
cd $FE

# install dependency
yarn install -f

# build code
yarn run build

# copy build files
rm -rf $TARGET_DIR
mkdir -p $TARGET_DIR
cp -rf dist/* $TARGET_DIR
cp frontend.conf $TARGET_DIR