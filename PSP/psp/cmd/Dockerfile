FROM centos:7

ENV WORKSPACE workspace
ENV YS_TOP /workspace

RUN mkdir -p /workspace
RUN mkdir -p /workspace/psp
RUN mkdir -p /workspace/psp/config
RUN chmod -R 777 /workspace


COPY pspd /workspace/psp/
COPY start.sh /workspace/psp/
COPY .env /workspace/psp/
COPY config /workspace/psp/config/

EXPOSE 32432

ENTRYPOINT ["/workspace/psp/start.sh"]
