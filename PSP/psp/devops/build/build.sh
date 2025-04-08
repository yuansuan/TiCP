#!/usr/bin/env bash
set -e

cd /workspace
mkdir -p dist || pwd

# Create bin directory
mkdir -p dist/bin || pwd

echo "backend start"
date

cd /workspace/cmd/ && make build && cd -

rm -rf /workspace/dist/schema
cp -rf devops/build/schema /workspace/dist
cp devops/build/bin/encrypt /workspace/dist/bin

cd /workspace/cmd/
cp pspd /workspace/dist/bin

# Copy psp service script
cp bin/psp.service.sh /workspace/dist/bin

cp .env /workspace/dist
sh sync_config.sh

cp -rf config /workspace/dist

# Remove some source codes
cd /workspace/dist/config/ && \rm -f config.go

cp /workspace/devops/build/rpm_resource/ysadmin /workspace/dist/bin
cd -

# Assign/remove the executable permission for these files
chmod +x /workspace/dist/bin/ysadmin
chmod +x /workspace/dist/bin/psp.service.sh
chmod +x /workspace/dist/bin/encrypt


echo "backend end"
date

