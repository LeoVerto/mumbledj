# Build environment for mumbledj - golang alpine container
FROM    golang:1.25.4-alpine@sha256:d3f0cf7723f3429e3f9ed846243970b20a2de7bae6a5b66fc5914e228d831bbb AS builder

# renovate: datasource=repology depName=alpine_3_22/opus-dev versioning=loose
ARG     OPUS_VERSION="1.5.2-r1"

ARG     branch=master

RUN     apk add --no-cache \
          ca-certificates \
          make \
          git \
          build-base \
          opus-dev=${OPUS_VERSION}

COPY    . $GOPATH/src/github.com/leoverto/mumbledj
WORKDIR $GOPATH/src/github.com/leoverto/mumbledj

RUN     make build-static install


# Export binary only from builder environment
FROM    alpine:3.22.2@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

# renovate: datasource=repology depName=alpine_3_22/ffmpeg versioning=loose
ARG     FFMPEG_VERSION="6.1.2-r2"
# renovate: datasource=repology depName=alpine_3_22/openssl versioning=loose
ARG     OPENSSL_VERSION="3.5.4-r0"
# renovate: datasource=repology depName=alpine_3_22/aria2 versioning=loose
ARG     ARIA2_VERSION="1.37.0-r1"
# renovate: datasource=repology depName=alpine_3_22/yt-dlp versioning=loose
ARG     YT_DLP_VERSION="2025.10.22-r0"

RUN     apk add --no-cache \
          ffmpeg=${FFMPEG_VERSION} \
          openssl=${OPENSSL_VERSION} \
          aria2=${ARIA2_VERSION} \
          yt-dlp=${YT_DLP_VERSION}

COPY    --from=builder /usr/local/bin/mumbledj /usr/local/bin/mumbledj

# Drop to user level privileges
RUN     addgroup -S mumbledj && adduser -S mumbledj -G mumbledj && chmod 750 /home/mumbledj
WORKDIR /home/mumbledj
USER    mumbledj
ENV     HOME=/home/mumbledj

RUN     mkdir -p .config/mumbledj && \
        mkdir -p .cache/mumbledj

ENTRYPOINT ["/usr/local/bin/mumbledj"]
