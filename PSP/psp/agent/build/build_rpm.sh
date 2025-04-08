#!/usr/bin/env bash
#
# Copyright (C) 2019 LambdaCal Inc.
#
# Build the rpm package
#
set -e
WORK_DIR=/workspace

cd /workspace/agent/ && make docker-build
chmod +x /workspace/agent/psp_agent

cd /workspace
mkdir -p dist

BRANCH=$1
VERSION=$2
# Get current time string as the rpm release tag
CURRENT_TIME=`date +%y%m%d%H%M`

rpmbuild  -bb \
          --define "_current_time ${CURRENT_TIME}" \
          --define "_version ${VERSION}" \
          ${WORK_DIR}/agent/build/rpm_resource/rpm_build.spec
# Move rpms to target dirs
cp $HOME/rpmbuild/RPMS/x86_64/psp-agent-${VERSION}-${CURRENT_TIME}.x86_64.rpm  ${WORK_DIR}/dist
