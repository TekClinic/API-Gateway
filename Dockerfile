# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Copy the source code
COPY . .

# Download git
RUN apk add --no-cache git

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway

FROM golang:1.22-alpine

COPY --from=builder /api-gateway /api-gateway

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
EXPOSE 8080

# Run
CMD ["/api-gateway"]