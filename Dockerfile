FROM golang:1.25.1-alpine

WORKDIR /app

RUN apk add --no-cache tzdata git \
    && go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8000