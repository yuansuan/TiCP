# Platform

sh devops/platform/build.sh

cd devops/platform && VERSION=master docker-compose -f docker-compose.yml up -d

## lookup docker image info
docker inspect idgen:master 

## Docker clean
docker system prune -f