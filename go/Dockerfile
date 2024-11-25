# Build stage
FROM golang:1.21-alpine AS builder

LABEL maintainer="Oyefule Oluwatayo <oyefuleoluwatayo@gmail.com>" \
  org.label-schema.name="Do Nginx" \
  org.label-schema.vendor="Oyefule Oluwatayo" \
  org.label-schema.schema-version="1.0"

WORKDIR /app
COPY . .

# Build the application
RUN go build -o /do-nginx

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /do-nginx /app/do-nginx

# Make sure to set permissions
RUN chmod +x /app/ssh-client

ENTRYPOINT ["/app/do-nginx"]
