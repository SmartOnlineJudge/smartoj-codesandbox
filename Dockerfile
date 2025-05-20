FROM golang:1.21 as builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /sandbox main.go

# 最终运行环境
FROM python:3.9-slim

WORKDIR /workspace

# 安装基础依赖
RUN apt-get update && apt-get install -y --no-install-recommends 
    gcc 
    libc-dev 
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# 复制编译好的二进制文件
COPY --from=builder /sandbox /sandbox

# 创建工作目录
RUN mkdir -p /tmp/sandbox

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["/sandbox"]