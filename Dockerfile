# Stage 1: Build static Go binary
FROM golang:1.26.6-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
RUN mkdir -p /app/data

COPY --from=builder /server /app/server
COPY --from=builder /app/docs /app/docs

EXPOSE 8000

ENV PORT=8000

CMD ["/app/server"]