package ftpd

import (
	"net"
	"strconv"
	"strings"
	"time"
)

var commandMap = map[string]func(*ftpConn, string){
	"LIST": commandLIST,
	"PASS": commandPASS,
	"PASV": commandPASV,
	"QUIT": commandQUIT,
	"SYST": commandSYST,
	"TYPE": commandTYPE,
	"USER": commandUSER,
}

func commandLIST(conn *ftpConn, param string) {
	conn.dataLock.Lock()
	defer conn.dataLock.Unlock()
	if conn.dataConn == nil {
		conn.writeMessage(425, "Error opening data socket")
		return
	}

	fileInfos, err := fsStore.ListDir(conn.curDir)
	if err != nil {
		conn.writeMessage(550, err.Error())
		return
	}

	conn.writeMessage(150, "Opening ASCII mode data connection for file list")
	content := ""
	for _, fileInfo := range fileInfos {
		content += fileInfo.Mode().String() +
			" 1 ftp ftp " +
			" " + strconv.Itoa(int(fileInfo.Size())) + " " +
			fileInfo.ModTime().Format(" Jan _2 15:04 ") +
			fileInfo.Name() + "\r\n"
	}
	conn.sendByteData([]byte(content))
}

func commandPASS(conn *ftpConn, param string) {
	conn.writeMessage(230, "Guest login ok")
}

func commandPASV(conn *ftpConn, param string) {
	if conn.dataConn != nil {
		conn.writeMessage(425, "Already connected")
		return
	}

	listener, err := net.ListenTCP("tcp", nil)
	if err != nil {
		conn.writeMessage(425, "Data connection failed")
		return
	}

	if listener.SetDeadline(time.Now().Add(10*time.Second)) != nil {
		conn.writeMessage(425, "Data connection failed")
		return
	}

	conn.dataLock.Lock()
	go func() {
		defer conn.dataLock.Unlock()
		conn.dataConn, _ = listener.Accept()
		listener.Close()
	}()

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
