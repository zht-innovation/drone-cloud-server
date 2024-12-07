FROM golang:1.23.4-alpine

WORKDIR /app/zht_cloud_server

COPY * .

RUN go install
