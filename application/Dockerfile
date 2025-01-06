FROM golang:1.23 AS build

RUN go install github.com/vektra/mockery/v2@latest

WORKDIR /go/src/app

COPY . .


RUN go get github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI


RUN rm -rf /go/src/app/internal/UI && \
    go build -o main cmd/main.go && \
    cp /go/src/app/internal/config/config.yaml config.yaml


# # Этап 1: Сборка приложения
# FROM golang:1.23 AS build

# # Устанавливаем рабочую директорию в контейнере
# WORKDIR /go/src/app

# # Копируем файлы проекта в контейнер
# COPY go.mod go.sum golangci.yaml ./
# COPY ./internal ./internal
# COPY ./cmd ./cmd
# COPY ./docs ./docs
# COPY ./static ./static

# # Загружаем зависимости и компилируем приложение
# RUN go mod download
# RUN go build -o main cmd/main.go

# # Этап 2: Минимальный образ
# FROM alpine:latest

# # Копируем только необходимое из этапа сборки
# COPY --from=build /go/src/app/main /go/src/app/main
# COPY --from=build /go/src/app/internal/config/config.yaml /go/src/app/config.yaml
# COPY --from=build /go/src/app/docs /go/src/app/docs
# COPY --from=build /go/src/app/static /go/src/app/static


# CMD ["/go/src/app/main"]

