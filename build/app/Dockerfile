FROM golang:1.26-alpine

WORKDIR /app
LABEL maintainer="pomapinigin@gmail.com"
COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../cmd/ ./cmd
COPY ../../internal/ ./internal
COPY ../../migrations/ ./migrations
COPY ../../config/ ./config
COPY ../../../pkg ./pkg

COPY ../../.env/ ./.env

RUN go build -o main ./cmd/app
EXPOSE 8080

CMD ["./main"]