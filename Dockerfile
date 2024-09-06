# Build environment for mumbledj - golang alpine container
FROM    golang:1.23.1-alpine@sha256:436e2d978524b15498b98faa367553ba6c3655671226f500c72ceb7afb2ef0b1 AS builder

# renovate: datasource=repology depName=alpine_3_20/opus-dev versioning=loose
ARG     OPUS_VERSION="1.5.2-r0"

ARG     branch=master
ENV     GO111MODULE=on

RUN     apk add --no-cache \
          ca-certificates \
          make \
          git \
          build-base \
          opus-dev=${OPUS_VERSION}

COPY    . $GOPATH/src/go.reik.pl/mumbledj

# add assets, which will be bundled with binary
WORKDIR $GOPATH/src/go.reik.pl/mumbledj
COPY    assets assets
RUN     make && make install


# Export binary only from builder environment
FROM    alpine:3.20.3@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d

# renovate: datasource=repology depName=alpine_3_20/ffmpeg versioning=loose
ARG     FFMPEG_VERSION="6.1.1-r8"
# renovate: datasource=repology depName=alpine_3_20/openssl versioning=loose
ARG     OPENSSL_VERSION="3.3.2-r0"
# renovate: datasource=repology depName=alpine_3_20/aria2 versioning=loose
ARG     ARIA2_VERSION="1.37.0-r0"
# renovate: datasource=repology depName=alpine_3_20/yt-dlp versioning=loose
ARG     YT_DLP_VERSION="2024.08.06-r0"

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

RUN     mkdir -p .config/mumbledj && \
        mkdir -p .cache/mumbledj

ENTRYPOINT ["/usr/local/bin/mumbledj"]
