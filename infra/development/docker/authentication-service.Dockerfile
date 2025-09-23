FROM golang:alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/services/authentication-service ./services/authentication-service
COPY backend/shared ./shared

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/authentication-service ./services/authentication-service/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/authentication-service /app/build/authentication-service
COPY --from=builder /app/shared /app/shared

ENTRYPOINT ["/app/build/authentication-service"]