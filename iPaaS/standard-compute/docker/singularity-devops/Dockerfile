FROM registry.intern.yuansuan.cn/euclid/ecp-build-go-bin as builder

WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN BUILD_ENV=dev make devops


FROM centos:7

COPY --from=builder /workspace/devops /usr/local/bin/devops
RUN yum install -y https://repo.ius.io/ius-release-el7.rpm \
    && yum install -y singularity
