FROM golang:1.26.1-alpine3.23 AS builder
# Set the working directory inside the container
WORKDIR /app

# Copy and configure the Go module
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go application with optimizations for smaller binary size
RUN go build -o groxy .

FROM alpine:3.23
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appuser \
 && adduser -S -G appuser -H -s /sbin/nologin appuser
 # Copy build
COPY --from=builder --chown=appuser:appuser /app/groxy /app/groxy

#Copy configuration file
COPY --from=builder --chown=appuser:appuser /app/config.yaml /app/config.yaml
USER appuser

# Expose the port that groxy listens on (default is 80)
EXPOSE 80
ENTRYPOINT ["/app/groxy"]