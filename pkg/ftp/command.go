package ftp

import (
	"net"
	"strconv"
	"strings"
)

var commandMap = map[string]func(*ftpConn, string){
	"PASS": commandPASS,
	"PASV": commandPASV,
	"QUIT": commandQUIT,
	"SYST": commandSYST,
	"TYPE": commandTYPE,
	"USER": commandUSER,
}

func commandPASS(conn *ftpConn, param string) {
	conn.writeMessage(230, "Guest login ok")
}

func commandPASV(conn *ftpConn, param string) {
	if conn.dataConn != nil {
		conn.writeMessage(425, "Already connected")
		return
	}

	listener, err := net.Listen("tcp", "")
	if err != nil {
		conn.writeMessage(425, "Data connection failed")
		return
	}
	go func() { conn.dataConn, err = listener.Accept() }()

	host, _, _ := net.SplitHostPort(conn.ctrlConn.LocalAddr().String())
	_, port, _ := net.SplitHostPort(listener.Addr().String())

	hostFields := strings.Split(host, ".")
	p, _ := strconv.Atoi(port)
	p1 := strconv.Itoa(p / 256)
	p2 := strconv.Itoa(p % 256)
	target := "(" + hostFields[0] + "," + hostFields[1] + "," + hostFields[2] + "," + hostFields[3] + "," + p1 + "," + p2 + ")"

	conn.writeMessage(227, "Entering Passive Mode "+target)
}

func commandQUIT(conn *ftpConn, param string) {
	conn.writeMessage(221, "Goodbye")
}

func commandSYST(conn *ftpConn, param string) {
	conn.writeMessage(215, "UNIX Type: L8")
}

func commandTYPE(conn *ftpConn, param string) {
	switch param {
	case "I":
		conn.writeMessage(200, "Type set to I")
	default:
		conn.writeMessage(500, "Type not supported")
	}
}

func commandUSER(conn *ftpConn, param string) {
	conn.writeMessage(230, "Guest login ok")
}
