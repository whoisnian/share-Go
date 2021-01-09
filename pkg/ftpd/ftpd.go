package ftpd

import (
	"bufio"
	"log"
	"net"
)

func handleConn(conn *ftpConn) {
	conn.writeMessage(220, "share-GO")
	for {
		line, err := conn.ctrlR.ReadString('\n')
		if err != nil {
			break
		}
		log.Println(line)
		conn.receiveLine(line)
	}
}

// Start listens on the addr and then creates goroutine to handle each connection.
func Start(addr string) {
	log.Printf("FTP server started: <ftp://%s>\n", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConn(&ftpConn{
			conn,
			bufio.NewWriter(conn),
			bufio.NewReader(conn),
			nil,
		})
	}
}
