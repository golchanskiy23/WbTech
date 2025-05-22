FROM golang:1.24 AS builder
WORKDIR /usr/local/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apt update && apt install -y curl
RUN go build -o main ./cmd/app/main.go

EXPOSE 3333
RUN useradd -m appuser
USER appuser

CMD ["./main"]