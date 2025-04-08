#!/usr/bin/env bash
#
# Copyright (C) 2019 LambdaCal Inc.
#
# Build the rpm package
#
set -e
WORK_DIR=/workspace
sh $WORK_DIR/devops/build/build.sh 1 2 3
BRANCH=$1
VERSION=$2
# Get current time string as the rpm release tag
CURRENT_TIME=`date +%y%m%d%H%M`
echo ${VERSION}-${CURRENT_TIME} > ${WORK_DIR}/dist/config/version

# Unnecessary: dependency Python3.6, not compatible python2.7，will cause rpmbuild error
rm -rf /workspace/dist/api/dist/onpremise/node_modules/node-gyp

rpmbuild  -bb \
          --define "_current_time ${CURRENT_TIME}" \
          --define "_version ${VERSION}" \
          ${WORK_DIR}/devops/build/rpm_resource/rpm_build.spec
# Move rpms to target dirs
cp $HOME/rpmbuild/RPMS/x86_64/psp-${VERSION}-${CURRENT_TIME}.x86_64.rpm  ${WORK_DIR}/dist
