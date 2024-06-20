# Build environment for mumbledj - golang alpine container
FROM    golang:1.22.4-alpine@sha256:f46eb222fb7fcd330618499eddef228fdc1be9b7eb2490d37139797e7e33ca38 AS builder

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
FROM    alpine:3.20.0@sha256:77726ef6b57ddf65bb551896826ec38bc3e53f75cdde31354fbffb4f25238ebd

# renovate: datasource=repology depName=alpine_3_20/ffmpeg versioning=loose
ARG     FFMPEG_VERSION="6.1.1-r8"
# renovate: datasource=repology depName=alpine_3_20/openssl versioning=loose
ARG     OPENSSL_VERSION="3.3.1-r0"
# renovate: datasource=repology depName=alpine_3_20/aria2 versioning=loose
ARG     ARIA2_VERSION="1.37.0-r0"
# renovate: datasource=repology depName=alpine_3_20/yt-dlp versioning=loose
ARG     YT_DLP_VERSION="2024.04.09-r1"

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
