FROM golang:1.24 AS builder
WORKDIR /usr/local/app

COPY go.mod go.sum ./
RUN go mod download

RUN apt update && apt install -y curl
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6

COPY . .
RUN golangci-lint run ./... --
RUN go build -o main ./cmd/app/main.go

EXPOSE 3333
RUN useradd -m appuser
USER appuser

CMD ["./main"]