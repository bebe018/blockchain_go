FROM golang:1.25.2 AS builder
WORKDIR /app
COPY . .
RUN go build -o blockchain .

FROM debian:bookworm-slim
WORKDIR /root
ENV NODE_ID=3000
COPY --from=builder /app/blockchain .

# 讓 container 一直活著（PID 1 是 tail）
ENTRYPOINT ["sh", "-c", "tail -f /dev/null"]
