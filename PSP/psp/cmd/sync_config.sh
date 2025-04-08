#!/usr/bin/env bash

cp -rf ${ROOTPATH}/internal/app/config/template config/
cp -rf ${ROOTPATH}/internal/rbac/config/rbac_model.conf config/
cp -rf ${ROOTPATH}/internal/user/config/license.yaml config/
cp -rf ${ROOTPATH}/internal/notice/config/sysconfig.yaml config/