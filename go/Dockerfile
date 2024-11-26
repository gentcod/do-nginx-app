# Build stage
FROM golang:1.21-alpine AS builder

LABEL maintainer="Oyefule Oluwatayo <oyefuleoluwatayo@gmail.com>"

WORKDIR /app
COPY . .
RUN go build -o /do-nginx

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /do-nginx /app/do-nginx
COPY .env .
COPY /scripts /app/scripts

RUN chmod +x /app/do-nginx

ENTRYPOINT ["/app/do-nginx"]
