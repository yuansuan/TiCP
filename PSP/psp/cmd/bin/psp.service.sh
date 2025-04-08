#!/usr/bin/env bash
#
# Copyright (C) 2020 LambdaCal Inc.
#
# chkconfig: 2345 99 90
#
# description: Start or stop PSP services for psp system service
#

ACTION=$1

# Default action is start
if [[ -z ${ACTION} ]] ; then
    ACTION=start
fi

# Start the PSP services
source @YS_TOP@/psp/config/profile;@YS_TOP@/psp/bin/ysadmin ${ACTION} all

# Get the CentOS/RHEL version number
OS_VERION=`cat /etc/redhat-release|sed -r 's/.* ([0-9]+)\..*/\1/'`

if [[ "${OS_VERION}x" != "6x" && "$1" == "start" ]] ; then
    while true ;
    do
        sleep 1
    done
fi