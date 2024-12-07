FROM golang:1.23.4-alpine

WORKDIR /app/zht_cloud_server

COPY * .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go get
RUN go run main.go
