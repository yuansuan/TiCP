# 使用一个轻量级的基础镜像来运行应用程序
FROM alpine:latest

# 创建工作目录
WORKDIR /root/

# 接受项目名作为构建参数
ARG PROJECT_NAME=mysql2mongo

# 复制二进制文件到工作目录
COPY ${PROJECT_NAME} .

# 暴露应用程序端口
EXPOSE 8080

# 启动应用程序
CMD ["sleep", "infinity"]
