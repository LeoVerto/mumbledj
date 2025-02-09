# Build environment for mumbledj - golang alpine container
FROM    golang:1.23.6-alpine@sha256:2c49857f2295e89b23b28386e57e018a86620a8fede5003900f2d138ba9c4037 AS builder

# renovate: datasource=repology depName=alpine_3_21/opus-dev versioning=loose
ARG     OPUS_VERSION="1.5.2-r1"

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
FROM    alpine:3.21.2@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099

# renovate: datasource=repology depName=alpine_3_21/ffmpeg versioning=loose
ARG     FFMPEG_VERSION="6.1.2-r1"
# renovate: datasource=repology depName=alpine_3_21/openssl versioning=loose
ARG     OPENSSL_VERSION="3.3.2-r5"
# renovate: datasource=repology depName=alpine_3_21/aria2 versioning=loose
ARG     ARIA2_VERSION="1.37.0-r0"
# renovate: datasource=repology depName=alpine_3_21/yt-dlp versioning=loose
ARG     YT_DLP_VERSION="2025.01.26-r0"

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
