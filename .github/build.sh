#!/bin/bash -ex
#########################################################################
# File Name: build.sh
# Author: nian
# Blog: https://whoisnian.com
# Mail: zhuchangbao1998@gmail.com
# Created Time: 2023年02月25日 星期六 20时00分56秒
#########################################################################

export CGO_ENABLED=0
export BUILDTIME=$(date --iso-8601=seconds)
if [[ -z "${GITHUB_REF_NAME}" ]]; then
  export VERSION=$(git describe --tags || echo unknown)
else
  export VERSION=${GITHUB_REF_NAME}
fi

goBuild() {
  GOOS="$1" GOARCH="$2" go build -trimpath \
    -ldflags="-s -w \
    -X 'github.com/whoisnian/share-Go/internal/global.Version=${VERSION}' \
    -X 'github.com/whoisnian/share-Go/internal/global.BuildTime=${BUILDTIME}'" \
    -o "$3" .
}

if [[ "$1" == '.' ]]; then
  goBuild $(go env GOOS) $(go env GOARCH) share-Go
elif [[ "$1" == 'linux-amd64' ]]; then
  goBuild linux amd64 "share-Go-linux-amd64-${VERSION}"
elif [[ "$1" == 'linux-arm64' ]]; then
  goBuild linux arm64 "share-Go-linux-arm64-${VERSION}"
elif [[ "$1" == 'darwin-amd64' ]]; then
  goBuild darwin amd64 "share-Go-darwin-amd64-${VERSION}"
elif [[ "$1" == 'darwin-arm64' ]]; then
  goBuild darwin arm64 "share-Go-darwin-arm64-${VERSION}"
elif [[ "$1" == 'windows-amd64' ]]; then
  goBuild windows amd64 "share-Go-windows-amd64-${VERSION}.exe"
elif [[ "$1" == 'windows-arm64' ]]; then
  goBuild windows arm64 "share-Go-windows-arm64-${VERSION}.exe"
elif [[ "$1" == 'all' ]]; then
  goBuild linux amd64 "share-Go-linux-amd64-${VERSION}"
  goBuild linux arm64 "share-Go-linux-arm64-${VERSION}"
  goBuild darwin amd64 "share-Go-darwin-amd64-${VERSION}"
  goBuild darwin arm64 "share-Go-darwin-arm64-${VERSION}"
  goBuild windows amd64 "share-Go-windows-amd64-${VERSION}.exe"
  goBuild windows arm64 "share-Go-windows-arm64-${VERSION}.exe"
else
  echo "Unknown build target"
  exit 1
fi
