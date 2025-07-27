FROM golang:1.22-alpine AS builder

WORKDIR /codesandbox

ENV GOPROXY=https://mirrors.aliyun.com/goproxy,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o sandbox main.go

FROM python:3.11-slim

WORKDIR /codesandbox

COPY --from=builder /codesandbox /codesandbox

EXPOSE 8080

CMD ["/codesandbox/sandbox"]
