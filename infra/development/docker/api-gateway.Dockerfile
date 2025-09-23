FROM golang:alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/services/api-gateway ./services/api-gateway
COPY backend/shared ./shared

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/api-gateway ./services/api-gateway/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/api-gateway /app/build/api-gateway
COPY --from=builder /app/shared /app/shared

ENTRYPOINT ["/app/build/api-gateway"]