
# --------------------------
# BUILD STAGE
# --------------------------
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build source
COPY cmd/ cmd/
COPY internal/ internal/
RUN go build -o aws-mcp-gateway ./cmd/mcp-server

# --------------------------
# RUNTIME STAGE
# --------------------------
FROM alpine:3.18

# Add CA certs
RUN apk add --no-cache ca-certificates

# Copy binary from builder
WORKDIR /app
COPY --from=builder /app/aws-mcp-gateway .


ENTRYPOINT [ "./aws-mcp-gateway" ]
