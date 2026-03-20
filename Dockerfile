FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o market-analysis-mcp .

FROM alpine:latest

RUN apk add --no-cache chromium

WORKDIR /app
COPY --from=builder /build/market-analysis-mcp .

EXPOSE 8080

ENTRYPOINT ["./market-analysis-mcp", "--serve"]
