#!/bin/bash

while true; do
    if [ -e "/etc/ys-agent/agent_env" ]; then
        echo "/etc/ys-agent/agent_env exists."
        break
    else
        echo "/etc/ys-agent/agent_env does not exist. Sleeping for 1 second..."
        sleep 1
    fi
done

. /etc/ys-agent/agent_env

while [ -z "$SHARE_SERVER" ]; do
    echo "SHARE_SERVER is empty. Waiting for it to be set..."
    sleep 1
done


# 目前SHARE_MOUNT_PATHS中设置了缺省挂载信息
echo "SHARE_MOUNT_PATHS is not empty. It is set to: $SHARE_MOUNT_PATHS"
IFS=',' read -ra paths_array <<< "$SHARE_MOUNT_PATHS"
IFS=',' read -ra share_username_array <<< "$SHARE_USERNAME"
share_num=${#share_username_array[@]}
for ((i=0; i<share_num; i++)); do
    path_pair=${paths_array[i]}
    IFS='=' read -ra pair <<< "$path_pair"

    share_username="${share_username_array[i]}"
    sub_path="${pair[0]}"
    mountPoint="${pair[1]}"

    echo "mountSubPath: ${share_username}, mountPoint: ${mountPoint}"

    mkdir -p ${mountPoint}
    mount -t cifs //${SHARE_SERVER}/${share_username} ${mountPoint} -o username=${share_username},password=${SHARE_PASSWORD},uid=2001,gid=1001
    ln -s ${mountPoint} /home/ecpuser/Desktop/$(basename $mountPoint)
done
