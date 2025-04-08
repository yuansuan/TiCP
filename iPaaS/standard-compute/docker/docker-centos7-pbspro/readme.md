### build
```
docker build . -t pbspro:v19.1
```


### run

```
docker run --rm  -it --hostname pbspro pbspro:v19.1
```

### submit job

```
docker exec -it ${CONTAINER_NAME} bash
```

```
[root@pbspro script]# su yskj
[yskj@pbspro script]$ cd
[yskj@pbspro ~]$ ls
[yskj@pbspro ~]$ echo "echo 123" >> script.sh
[yskj@pbspro ~]$ qsub script.sh
0.pbspro
[yskj@pbspro ~]$ qstat
[yskj@pbspro ~]$ qstat -x 0
Job id            Name             User              Time Use S Queue
----------------  ---------------- ----------------  -------- - -----
0.pbspro          script.sh        yskj              00:00:00 F workq
```