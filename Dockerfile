# Copyright 2024 Deutsche Telekom IT GmbH
#
# SPDX-License-Identifier: Apache-2.0

FROM --platform=$BUILDPLATFORM golang:1.27.1-alpine AS build
ARG GOPROXY
ARG GONOSUMDB
ARG TARGETOS
ARG TARGETARCH
ENV GOPROXY=$GOPROXY
ENV GONOSUMDB=$GONOSUMDB

USER root

WORKDIR /build
COPY . .
RUN apk add --no-cache build-base
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s -extldflags=-static" -o cosmoparrot

FROM scratch
ARG VERSION
ARG REVISION
LABEL org.opencontainers.image.source="https://github.com/telekom/pubsub-horizon-cosmoparrot" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${REVISION}" \
      org.opencontainers.image.licenses="Apache-2.0"
COPY --from=build /build/cosmoparrot cosmoparrot

USER 65532:65532
ENTRYPOINT ["./cosmoparrot"]
