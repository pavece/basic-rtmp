FROM golang:1.23-bookworm AS builder 
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o rtmp-hls ./cmd/basic-rtmp

FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/rtmp-hls /app/rtmp-hls

# Set pipe buffer size limit
RUN echo "* soft pipe 8192" >> /etc/security/limits.conf && \
    echo "* hard pipe 8192" >> /etc/security/limits.conf

EXPOSE 1935
CMD ["sh", "-c", "ulimit -p 8192 && /app/rtmp-hls"]

