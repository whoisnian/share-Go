package ftpd

import (
	"net"

	"github.com/whoisnian/glb/logger"
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
		if conn.closed {
			break
		}
	}
}

// Start listens on the addr and then creates goroutine to handle each connection.
func Start(addr string, rootPath string) {
	var err error
	if fsStore, err = storage.New(rootPath); err != nil {
		logger.Fatal(err)
	}

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
