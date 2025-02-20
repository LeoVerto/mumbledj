# Build environment for mumbledj - golang alpine container
FROM    golang:1.24.0-alpine@sha256:2d40d4fc278dad38be0777d5e2a88a2c6dee51b0b29c97a764fc6c6a11ca893c AS builder

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
FROM    alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

# renovate: datasource=repology depName=alpine_3_21/ffmpeg versioning=loose
ARG     FFMPEG_VERSION="6.1.2-r1"
# renovate: datasource=repology depName=alpine_3_21/openssl versioning=loose
ARG     OPENSSL_VERSION="3.3.3-r0"
# renovate: datasource=repology depName=alpine_3_21/aria2 versioning=loose
ARG     ARIA2_VERSION="1.37.0-r0"
# renovate: datasource=repology depName=alpine_3_21/yt-dlp versioning=loose
ARG     YT_DLP_VERSION="2025.02.19-r0"

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
