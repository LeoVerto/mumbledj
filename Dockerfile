# Build environment for mumbledj - golang alpine container
FROM golang:1.24-alpine AS builder
ARG branch=master

RUN apk add --no-cache ca-certificates make git build-base opus-dev
COPY . $GOPATH/src/go.reik.pl/mumbledj

# add assets, which will be bundled with binary
WORKDIR $GOPATH/src/go.reik.pl/mumbledj
COPY assets assets
RUN make && make install


# Export binary only from builder environment
FROM alpine:latest
RUN apk add --no-cache ffmpeg openssl aria2 yt-dlp
COPY --from=builder /usr/local/bin/mumbledj /usr/local/bin/mumbledj

# Drop to user level privileges
RUN addgroup -S mumbledj && adduser -S mumbledj -G mumbledj && chmod 750 /home/mumbledj
WORKDIR /home/mumbledj
USER mumbledj
ENV HOME=/home/mumbledj
RUN mkdir -p .config/mumbledj && \
    mkdir -p .cache/mumbledj

ENTRYPOINT ["/usr/local/bin/mumbledj"]
