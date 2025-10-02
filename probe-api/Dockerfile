# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o probe-api ./cmd/probe-api

# Run stage
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /app/probe-api .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/probe-api"]
