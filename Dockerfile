# Stage 1: Build
FROM golang:1.20.5-alpine3.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Create a non-root user and group
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Copy the binary and change ownership
COPY --from=builder /app/main /app/main

# Copy the templates directory into the container
COPY templates/ /app/templates/

RUN chown appuser:appgroup /app/main

# Switch to the non-root user
USER appuser

EXPOSE 8080

CMD ["/app/main"]