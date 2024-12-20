FROM golang:1.23 AS build

RUN go install github.com/vektra/mockery/v2@latest

WORKDIR /go/src/app

COPY . .

# RUN go get github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI

RUN rm -rf /go/src/app/internal/UI && \
    go build -o main cmd/main.go && \
    cp /go/src/app/internal/config/config.yaml config.yaml
