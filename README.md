# share-Go
[![Release Status](https://github.com/whoisnian/share-Go/actions/workflows/release.yml/badge.svg)](https://github.com/whoisnian/share-Go/actions/workflows/release.yml)

## Usage
Download the latest binary from [Release Page](https://github.com/whoisnian/share-Go/releases) according to your operating system and architecture. Alternatively, docker container is also supported.
### Run binary directly
```sh
mkdir ./uploads
./share-Go -log nano -l 0.0.0.0:9000 -p ./uploads
```
### With linux chroot
```sh
mkdir -p ./share/uploads && sudo chown 65534:65534 ./share/uploads

# initialize chroot environment with alpine minirootfs
wget https://dl-cdn.alpinelinux.org/alpine/v3.21/releases/x86_64/alpine-minirootfs-3.21.3-x86_64.tar.gz
tar -xvf alpine-minirootfs-3.21.3-x86_64.tar.gz -C ./share && rm alpine-minirootfs-3.21.3-x86_64.tar.gz
echo 'nameserver 223.5.5.5' > ./share/etc/resolv.conf

# move share-Go binary into ./share and run
sudo chroot --userspec=65534:65534 ./share /share-Go -log text -l 0.0.0.0:9000 -p /uploads
```
example `/etc/systemd/system/share-Go.service`:
```
[Unit]
Description=share-Go Service
After=network-online.target

[Service]
Type=simple
User=root
Restart=always
RestartSec=5s
ExecStart=/usr/bin/chroot --userspec=65534:65534 /root/share /share-Go -log text -l 0.0.0.0:9000 -p /uploads

[Install]
WantedBy=multi-user.target
```
### With docker container
```sh
mkdir ./uploads
docker run -d \
  --name share-go \
  -e CFG_LOG_FMT=json \
  -e CFG_LISTEN_ADDR=:9000 \
  -e CFG_ROOT_PATH=/uploads \
  -p 9000:9000 \
  -v $(pwd)/uploads:/uploads \
  ghcr.io/whoisnian/share-go:v0.0.11
```

## Development
* start backend service:
  ```sh
  mkdir ./uploads
  go run ./main.go -log nano -d # manually rerun after modifying the golang code
  ```
* start frontend dev server:
  ```sh
  cd fe
  npm install
  npm run start # live reloading for the javascript code
  ```
* visit http://127.0.0.1:9100 in your web browser.
