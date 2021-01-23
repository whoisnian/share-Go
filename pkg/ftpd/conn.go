package ftpd

import (
	"bufio"
	"io"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/whoisnian/share-Go/pkg/logger"
)

type ftpConn struct {
	ctrlConn net.Conn
	ctrlW    *bufio.Writer
	ctrlR    *bufio.Reader
	dataConn net.Conn
	dataLock sync.Mutex
	curDir   string
}

func newftpConn(conn net.Conn) *ftpConn {
	return &ftpConn{
		ctrlConn: conn,
		ctrlW:    bufio.NewWriter(conn),
		ctrlR:    bufio.NewReader(conn),
		dataConn: nil,
		curDir:   "/",
	}
}

func (conn *ftpConn) writeMessage(code int, message string) {
	content := strconv.Itoa(code) + " " + message + "\r\n"
	logger.Debug("<-- ", content)
	conn.ctrlW.WriteString(content)
	conn.ctrlW.Flush()
}

func (conn *ftpConn) sendByteData(data []byte) error {
	defer func() {
		conn.dataConn.Close()
		conn.dataConn = nil
	}()

	if _, err := conn.dataConn.Write(data); err != nil {
		return err
	}
	conn.writeMessage(226, "Data transmission OK")
	return nil
}

func (conn *ftpConn) sendStreamData(reader io.ReadCloser) error {
	defer func() {
		conn.dataConn.Close()
		conn.dataConn = nil
	}()

	if _, err := io.Copy(conn.dataConn, reader); err != nil {
		return err
	}
	conn.writeMessage(226, "Data transmission OK")
	return nil
}

func (conn *ftpConn) writeStreamData(writer io.WriteCloser) error {
	defer func() {
		conn.dataConn.Close()
		conn.dataConn = nil
	}()

	if _, err := io.Copy(writer, conn.dataConn); err != nil {
		return err
	}
	conn.writeMessage(226, "Data transmission OK")
	return nil
}

func (conn *ftpConn) receiveLine(line string) {
	logger.Debug("--> ", line)
	command, param := conn.parseLine(line)
	if commandFunc, ok := commandMap[strings.ToUpper(command)]; ok {
		commandFunc(conn, param)
	} else {
		conn.writeMessage(502, "Command not implemented")
		return
	}
}

func (conn *ftpConn) parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}

func (conn *ftpConn) buildPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(conn.curDir, path)
}
