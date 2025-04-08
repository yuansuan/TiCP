FROM centos:7

MAINTAINER cchu <cchu@yuansuan.cn>

RUN yum -y install epel-release \
    && yum -y install gcc file tree mysql \
    && yum -y install openssh-server openssh-clients sshpass pam-devel \
    && yum -y install unzip \
    && yum -y install git \
    && yum -y install sudo \
    && yum -y install wget \
    && yum clean all

RUN yum install -y expat libedit postgresql-server postgresql-contrib python \
sendmail sudo tcl tk libical

#  install pbs
RUN mkdir -p /usr/local/pbs
RUN wget https://static-1252829527.cos.ap-shanghai.myqcloud.com/pbspro/pbspro-debuginfo-19.1.2-0.x86_64.rpm -O /usr/local/pbs/pbspro-debuginfo-19.1.2-0.x86_64.rpm
RUN wget https://static-1252829527.cos.ap-shanghai.myqcloud.com/pbspro/pbspro-server-19.1.2-0.x86_64.rpm -O /usr/local/pbs/pbspro-server-19.1.2-0.x86_64.rpm
WORKDIR /usr/local/pbs
RUN yum install -y pbspro-server-19.1.2-0.x86_64.rpm pbspro-debuginfo-19.1.2-0.x86_64.rpm
ADD resource/script /workspace


WORKDIR /workspace
ENTRYPOINT ["bash", "entrypoint.sh"]

