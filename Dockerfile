FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev git protobuf

ENV GOBIN=/usr/local/bin

RUN go install github.com/bufbuild/buf/cmd/buf@v1.28.1 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY proto/ proto/
COPY buf.gen.yaml buf.yaml ./
RUN buf generate

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/smart-hub ./cmd/api/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/smart-hub .

EXPOSE 50051

ENV SERVICE_NAME=smart-hub \
    SERVICE_PORT=50051 \
    SERVICE_ENV=prod \
    LOG_LEVEL=INFO

CMD ["./smart-hub"]