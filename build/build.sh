#!/bin/bash -e
#########################################################################
# File Name: build.sh
# Author: nian
# Blog: https://whoisnian.com
# Mail: zhuchangbao1998@gmail.com
# Created Time: 2024年02月10日 星期六 20时10分28秒
#########################################################################

SCRIPT_DIR=$(dirname "$0")
SOURCE_DIR="$SCRIPT_DIR/.."
OUTPUT_DIR="$SOURCE_DIR/output"

MODULE_NAME=$(go mod edit -fmt -print | grep -Po '(?<=^module ).*$')
APP_NAME=$(basename "$MODULE_NAME")
BUILDTIME=$(date --iso-8601=seconds)
if [[ -z "$GITHUB_REF_NAME" ]]; then
  VERSION=$(git describe --tags 2> /dev/null || echo unknown)
else
  VERSION=$GITHUB_REF_NAME
fi

goBuild() {
  CGO_ENABLED=0 GOOS="$1" GOARCH="$2" go build \
    -trimpath -ldflags="-s -w \
    -X '${MODULE_NAME}/global.ModName=${MODULE_NAME}' \
    -X '${MODULE_NAME}/global.AppName=${APP_NAME}' \
    -X '${MODULE_NAME}/global.Version=${VERSION}' \
    -X '${MODULE_NAME}/global.BuildTime=${BUILDTIME}'" \
    -o "$OUTPUT_DIR"/"$3" "$SOURCE_DIR"
}

if [[ "$1" == '.' ]]; then
  goBuild $(go env GOOS) $(go env GOARCH) "$APP_NAME"
elif [[ "$1" == 'all' ]]; then
  goBuild linux amd64 "${APP_NAME}-linux-amd64-${VERSION}"
  goBuild linux arm64 "${APP_NAME}-linux-arm64-${VERSION}"
  goBuild darwin amd64 "${APP_NAME}-darwin-amd64-${VERSION}"
  goBuild darwin arm64 "${APP_NAME}-darwin-arm64-${VERSION}"
  goBuild windows amd64 "${APP_NAME}-windows-amd64-${VERSION}.exe"
  goBuild windows arm64 "${APP_NAME}-windows-arm64-${VERSION}.exe"
elif [[ "$#" == 2 ]]; then
  goBuild "$1" "$2" "${APP_NAME}-$1-$2-${VERSION}"
else
  cat << EOF
Usage:
  $(basename $0) .            # build for current platform $(go env GOOS)-$(go env GOARCH)
  $(basename $0) all          # build for all supported platforms
  $(basename $0) darwin amd64 # build for specified platform

Supported platforms:
  linux-amd64
  linux-arm64
  darwin-amd64
  darwin-arm64
  windows-amd64
  windows-arm64
EOF
  exit 1
fi
