# syntax=docker.io/docker/dockerfile:1.14
# https://docs.docker.com/build/dockerfile/frontend/

FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.24-alpine AS build

ARG MODULE_NAME
ARG APP_NAME
ARG VERSION
ARG BUILDTIME

# https://docs.docker.com/build/guide/mounts/
WORKDIR /app
RUN --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x

# https://docs.docker.com/build/building/variables/#scoping
ARG TARGETOS
ARG TARGETARCH

# https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build/,id=build-$TARGETARCH \
    --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w \
    -X '${MODULE_NAME}/global.ModName=${MODULE_NAME}' \
    -X '${MODULE_NAME}/global.AppName=${APP_NAME}' \
    -X '${MODULE_NAME}/global.Version=${VERSION}' \
    -X '${MODULE_NAME}/global.BuildTime=${BUILDTIME}'" \
    -o main .

# RedundantTargetPlatform: setting --platform=$TARGETPLATFORM in FROM is the default behavior
FROM gcr.io/distroless/static-debian12:latest
COPY --from=build /app/main /main
ENTRYPOINT ["/main"]
