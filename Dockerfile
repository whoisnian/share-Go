# syntax=docker.io/docker/dockerfile:1.6

FROM golang:1.21-alpine AS build

ARG VERSION
ARG BUILDTIME

WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build,id=build-$ARCH \
    --mount=type=cache,target=/go/pkg/mod \
  CGO_ENABLED=0 go build -trimpath \
  -ldflags="-s -w \
  -X 'github.com/whoisnian/share-Go/internal/global.Version=${VERSION}' \
  -X 'github.com/whoisnian/share-Go/internal/global.BuildTime=${BUILDTIME}'" \
  -o share-Go .

FROM gcr.io/distroless/static-debian12:latest
COPY --from=build /app/share-Go /
ENTRYPOINT ["/share-Go"]
