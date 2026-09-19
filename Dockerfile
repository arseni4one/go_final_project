FROM golang:1.26-alpine AS builder
LABEL authors="arseni4"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo ./todo-app
COPY --from=builder /app/web ./web

ENV TODO_PORT=7540
ENV TODO_PASSWORD="12345"
ENV TODO_DBFILE="./scheduler.db"



ENTRYPOINT ["./todo-app"]