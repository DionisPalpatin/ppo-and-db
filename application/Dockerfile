FROM golang:1.23 AS build-stage

WORKDIR /notebook-app-server

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server.exe cmd/main.go


FROM alpine

WORKDIR /

COPY --from=build-stage /notebook-app-server/server.exe /servers/notebook-app/application/server.exe
COPY ./static /servers/notebook-app/static
COPY ./docs /servers/notebook-app/docs
