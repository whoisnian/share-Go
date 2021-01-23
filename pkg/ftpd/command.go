package ftpd

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/whoisnian/share-Go/pkg/util"
)

var commandMap = map[string]func(*ftpConn, string){
	"CWD": commandCWD,
	"LIST": commandLIST,
	"PASS": commandPASS,
	"PASV": commandPASV,
	"PWD":  commandPWD,
	"QUIT": commandQUIT,
	"RETR": commandRETR,
	"STOR": commandSTOR,
	"SYST": commandSYST,
	"TYPE": commandTYPE,
	"USER": commandUSER,
}

func commandCWD(conn *ftpConn, param string) {
	path := strings.TrimSpace(param)
	if len(path) < 1 {
		path = "/"
	} else if filepath.IsAbs(path) {
		// noop
	} else {
		path = filepath.Join(conn.curDir, path)
	}

	if !fsStore.IsDir(path) {
		conn.writeMessage(550, "No such directory")
		return
	}

	conn.curDir = path
	conn.writeMessage(250, "Change working directory successfully")
}

func commandLIST(conn *ftpConn, param string) {
	conn.dataLock.Lock()
	defer conn.dataLock.Unlock()
	if conn.dataConn == nil {
		conn.writeMessage(425, "Error opening data socket")
		return
	}

	i := 0
	for ; i < len(param); i++ {
		if (i == 0 || util.IsSpace(param[i-1])) && !util.IsSpace(param[i]) && param[i] != '-' {
			break
		}
	}

	fileInfos, err := fsStore.ListDir(conn.buildPath(param[i:]))
	if err != nil {
		conn.writeMessage(550, err.Error())
		return
	}

	content := ""
	for _, fileInfo := range fileInfos {
		content += fileInfo.Mode().String() +
			" 1 ftp ftp " +
			" " + strconv.Itoa(int(fileInfo.Size())) + " " +
			fileInfo.ModTime().Format(" Jan _2 15:04 ") +
			fileInfo.Name() + "\r\n"
	}
	conn.writeMessage(150, "Opening ASCII mode data connection for file list")
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

func commandPWD(conn *ftpConn, param string) {
	conn.writeMessage(257, "\""+conn.curDir+"\" is current directory")
}

func commandQUIT(conn *ftpConn, param string) {
	conn.writeMessage(221, "Goodbye")
}

func commandRETR(conn *ftpConn, param string) {
	conn.dataLock.Lock()
	defer conn.dataLock.Unlock()
	if conn.dataConn == nil {
		conn.writeMessage(425, "Error opening data socket")
		return
	}

	file, err := fsStore.GetFile(conn.buildPath(param))
	if err != nil {
		if os.IsNotExist(err) {
			conn.writeMessage(550, "File does not exist")
		} else if os.IsPermission(err) {
			conn.writeMessage(550, "File permission is denied")
		} else {
			conn.writeMessage(550, "Unexpected error")
		}
		return
	}

	conn.writeMessage(150, "Sending file")
	conn.sendStreamData(file)
}

func commandSTOR(conn *ftpConn, param string) {
	conn.dataLock.Lock()
	defer conn.dataLock.Unlock()
	if conn.dataConn == nil {
		conn.writeMessage(425, "Error opening data socket")
		return
	}

	file, err := fsStore.CreateFile(conn.buildPath(param))
	if err != nil {
		if os.IsPermission(err) {
			conn.writeMessage(550, "File permission is denied")
		} else {
			conn.writeMessage(550, "Unexpected error")
		}
		return
	}

	conn.writeMessage(150, "Receiving file")
	conn.writeStreamData(file)
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
