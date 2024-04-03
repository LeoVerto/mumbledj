# Build environment for mumbledj - golang alpine container
FROM    golang:1.22.2-alpine@sha256:cdc86d9f363e8786845bea2040312b4efa321b828acdeb26f393faa864d887b0 AS builder

# renovate: datasource=repology depName=alpine_3_19/opus-dev versioning=loose
ARG     OPUS_VERSION="1.4-r0"

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
FROM    alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

# renovate: datasource=repology depName=alpine_3_19/ffmpeg versioning=loose
ARG     FFMPEG_VERSION="6.1.1-r0"
# renovate: datasource=repology depName=alpine_3_19/openssl versioning=loose
ARG     OPENSSL_VERSION="3.1.4-r5"
# renovate: datasource=repology depName=alpine_3_19/aria2 versioning=loose
ARG     ARIA2_VERSION="1.37.0-r0"
# renovate: datasource=repology depName=alpine_3_19/yt-dlp versioning=loose
ARG     YT_DLP_VERSION="2023.11.16-r0"

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
