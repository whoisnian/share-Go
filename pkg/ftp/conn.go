package ftp

import (
	"bufio"
	"net"
	"strconv"
	"strings"
)

type ftpConn struct {
	ctrlConn net.Conn
	ctrlW    *bufio.Writer
	ctrlR    *bufio.Reader
	dataConn net.Conn
}

func (conn *ftpConn) writeMessage(code int, message string) {
	content := strconv.Itoa(code) + " " + message + "\r\n"
	conn.ctrlW.WriteString(content)
	conn.ctrlW.Flush()
}

func (conn *ftpConn) receiveLine(line string) {
	command, param := conn.parseLine(line)
	commandFunc := commandMap[strings.ToUpper(command)]
	if commandFunc == nil {
		conn.writeMessage(502, "Command not implemented")
		return
	}
	commandFunc(conn, param)
}

func (conn *ftpConn) parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}
