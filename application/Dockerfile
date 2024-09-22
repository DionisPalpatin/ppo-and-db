FROM golang:latest

WORKDIR /notebookapp

COPY . .

# RUN go get github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI

RUN rm -rf /notebookapp/internal/UI/GUI && \
    go build -o main.exe cmd/main.go

CMD ["/notebookapp/main.exe"]