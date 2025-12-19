FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# --- run stage ---
FROM gcr.io/distroless/base-debian12

WORKDIR /
COPY --from=builder /app/api /api

EXPOSE 8081
ENTRYPOINT ["/api"]
