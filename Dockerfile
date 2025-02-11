FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o cdn ./cmd/cdn

FROM alpine
COPY --from=builder /app/cdn /app/
CMD ["/app/cdn"]