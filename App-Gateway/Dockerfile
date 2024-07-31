
FROM golang:1.20.5-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOPATH=/go
ENV GOCACHE=/go-build


RUN go install github.com/cosmtrek/air@v1.44.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
