FROM golang:1.25.2 AS builder
WORKDIR /app
COPY . .
RUN go build -o blockchain .

FROM debian:bookworm-slim
WORKDIR /root
ENV NODE_ID=3000
COPY --from=builder /app/blockchain .
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

EXPOSE 2112
ENTRYPOINT ["sh", "-c", "./blockchain & tail -f /dev/null"]

