# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS base

# Install dependencies only when needed
FROM base AS deps

# Download git
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod .
COPY go.sum .

# Build argument for GitHub token
ARG GITHUB_TOKEN

# Set up git configuration to use token for private repo
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

# Download Go modules
RUN go mod download

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps $GOPATH $GOPATH

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway

# Production image, copy all the files and run
FROM golang:1.22-alpine

COPY --from=builder /api-gateway /api-gateway

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
EXPOSE 8080

# Run
CMD ["/api-gateway"]