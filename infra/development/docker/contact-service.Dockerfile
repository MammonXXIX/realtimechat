FROM golang:alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/services/contact-service ./services/contact-service
COPY backend/shared ./shared

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/contact-service ./services/contact-service/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/contact-service /app/build/contact-service
COPY --from=builder /app/shared /app/shared

ENTRYPOINT ["/app/build/contact-service"]
