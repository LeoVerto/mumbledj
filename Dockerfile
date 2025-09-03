# Build environment for mumbledj - golang alpine container
FROM    golang:1.25.1-alpine@sha256:cb0b8e92b8b63b1ba16ce78d926fcba5d4f9ff241855115568006affc3ae6557 AS builder

# renovate: datasource=repology depName=alpine_3_22/opus-dev versioning=loose
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
FROM    alpine:3.22.1@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

# renovate: datasource=repology depName=alpine_3_22/ffmpeg versioning=loose
ARG     FFMPEG_VERSION="6.1.2-r2"
# renovate: datasource=repology depName=alpine_3_22/openssl versioning=loose
ARG     OPENSSL_VERSION="3.5.2-r0"
# renovate: datasource=repology depName=alpine_3_22/aria2 versioning=loose
ARG     ARIA2_VERSION="1.37.0-r1"
# renovate: datasource=repology depName=alpine_3_22/yt-dlp versioning=loose
ARG     YT_DLP_VERSION="2025.08.27-r0"

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
