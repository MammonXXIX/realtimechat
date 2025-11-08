FROM golang:alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/services/chat-service ./services/chat-service
COPY backend/shared ./shared

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/chat-service ./services/chat-service/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/chat-service /app/build/chat-service
COPY --from=builder /app/shared /app/shared

ENTRYPOINT ["/app/build/chat-service"]
