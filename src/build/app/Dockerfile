FROM golang:1.26-alpine

WORKDIR /app
LABEL maintainer="pomapinigin@gmail.com"
COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../cmd/ ./cmd
COPY ../../internal/ ./internal
COPY ../../config/ ./config
COPY ../../.env/ ./.env

RUN go build -o main ./cmd/app
EXPOSE 8080
CMD ["./main"]