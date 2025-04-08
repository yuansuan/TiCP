FROM harbor.yuansuan.cn/ticp/rockylinux:9

RUN useradd slurm && useradd -m yuansuan

ARG MODULENAME

WORKDIR /workspace/
ENV LOG_LVL="info"

ADD $MODULENAME /workspace/$MODULENAME/

ENV COMMAND=/workspace/$MODULENAME/$MODULENAME
WORKDIR /workspace/$MODULENAME

ENTRYPOINT ["/bin/sh", "-c", "$COMMAND"]
