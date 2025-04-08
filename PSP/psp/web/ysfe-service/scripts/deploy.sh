npm run doc:build

ssh root@10.0.1.25 rm -rf ysfe-service-doc
ssh root@10.0.1.25 mkdir ysfe-service-doc
scp -r doc/build/* root@10.0.1.25:/root/ysfe-service-doc/