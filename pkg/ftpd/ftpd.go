package ftpd

import (
	"net"

	"github.com/whoisnian/share-Go/pkg/logger"
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

	logger.Info("Service ftpd started: <ftp://", addr, ">")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Panic(err)
		}
		go handleConn(newftpConn(conn))
	}
}
