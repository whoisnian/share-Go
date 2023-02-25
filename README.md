# share-Go
[![Release Status](https://github.com/whoisnian/share-Go/workflows/release/badge.svg)](https://github.com/whoisnian/share-Go/actions?query=workflow%3Arelease)

## Usage
Download the latest binary from [Release Page](https://github.com/whoisnian/share-Go/releases) according to your operating system and architecture.
### Run in LAN
```sh
mkdir ./uploads
./share-Go -l 0.0.0.0:9000 -p ./uploads
```
### With linux chroot
```sh
mkdir -p ./share/uploads
# move share-Go binary into ./share, like:
# share/
#   ├── share-Go
#   └── uploads/
sudo chroot --userspec=$(id -u):$(id -g) ./share ./share-Go -l 0.0.0.0:9000 -p ./uploads
```
### As systemd service
`/etc/systemd/system/share-Go.service`
```
[Unit]
Description=share-Go Service
After=network-online.target

[Service]
Type=simple
User=root
Restart=always
RestartSec=5s
ExecStart=/usr/bin/chroot --userspec=nobody:nobody /root/share ./share-Go -l 0.0.0.0:9000 -p ./uploads

[Install]
WantedBy=multi-user.target
```

## Development
* start backend service:
  ```sh
  mkdir ./uploads
  go run ./main.go # manually rerun after modifying the golang code
  ```
* start frontend dev server:
  ```sh
  cd fe
  npm install
  npm run start # live reloading for the javascript code
  ```
* visit http://127.0.0.1:9100 in your web browser.
