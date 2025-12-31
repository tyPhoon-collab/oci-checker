# Build stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
# No go.sum yet, so we generate it
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o oci-checker .

# Run stage
FROM alpine:latest

WORKDIR /app

# OCI SDK requires CA certificates
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/oci-checker .

CMD ["./oci-checker"]
