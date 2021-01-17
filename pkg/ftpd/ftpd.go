package ftpd

import (
	"log"
	"net"

	"github.com/whoisnian/share-Go/pkg/storage"
)

var fsStore *storage.Store

func handleConn(conn *ftpConn) {
	conn.writeMessage(220, "share-GO")
	for {
		line, err := conn.ctrlR.ReadString('\n')
		if err != nil {
			break
		}
		conn.receiveLine(line)
	}
}

// Start listens on the addr and then creates goroutine to handle each connection.
func Start(addr string, rootPath string) {
	fsStore = storage.New(rootPath)

	log.Printf("Service ftpd started: <ftp://%s>\n", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConn(newftpConn(conn))
	}
}
