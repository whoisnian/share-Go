package ftpd

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/whoisnian/share-Go/pkg/util"
)

var commandMap = map[string]func(*ftpConn, string){
	"CWD":  commandCWD,
	"DELE": commandDELE,
	"FEAT": commandFEAT,
	"LIST": commandLIST,
	"MDTM": commandMDTM,
	"MKD":  commandMKD,
	"MLSD": commandMLSD,
	"MLST": commandMLST,
	"NOOP": commandNOOP,
	"OPTS": commandOPTS,
	"PASS": commandPASS,
	"PASV": commandPASV,
	"PWD":  commandPWD,
	"QUIT": commandQUIT,
	"RETR": commandRETR,
	"RMD":  commandRMD,
	"SIZE": commandSIZE,
	"STOR": commandSTOR,
	"SYST": commandSYST,
	"TYPE": commandTYPE,
	"USER": commandUSER,
}

func commandCWD(conn *ftpConn, param string) {
	path := strings.TrimSpace(param)
	if len(path) < 1 {
		path = "/"
	} else {
		path = conn.buildPath(path)
	}

	if !fsStore.IsDir(path) {
		conn.writeMessage(550, "No such directory")
		return
	}

	conn.curDir = path
	conn.writeMessage(250, "Change working directory successfully")
}

func commandDELE(conn *ftpConn, param string) {
	path := strings.TrimSpace(param)
	if len(path) < 1 {
		conn.writeMessage(500, "Syntax error")
		return
	}
	path = conn.buildPath(path)

	if !fsStore.IsFile(path) {
		conn.writeMessage(550, "No such file")
		return
	}

	err := fsStore.Delete(path)
	if err != nil {
		conn.writeMessage(550, err.Error())
		return
	}

	conn.writeMessage(250, "Delete file successfully")
}

func commandFEAT(conn *ftpConn, param string) {
	extendCommands := []string{
		"MDTM",
		"MLST Type;Perm;Modify;Size;",
		"SIZE",
		"UTF8",
	}

	content := "Supported extensions:\r\n"
	for _, cmd := range extendCommands {
		content += " " + cmd + "\r\n"
	}
	conn.writeMessageMultiline(211, content)
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

	fileInfo, err := fsStore.FileInfo(conn.buildPath(param[i:]))
	if err != nil {
		conn.writeMessage(550, "File or directory not available")
		return
	}

	content := ""
	if fileInfo.Mode().IsRegular() {
		content = fileInfo.Mode().String() +
			" 1 ftp ftp " +
			" " + strconv.Itoa(int(fileInfo.Size())) + " " +
			fileInfo.ModTime().Format(" Jan _2 15:04 ") +
			fileInfo.Name() + "\r\n"
	} else {
		infos, err := fsStore.ListDir(conn.buildPath(param[i:]))
		if err != nil {
			conn.writeMessage(550, "Directory not available")
			return
		}

		for _, info := range infos {
			content += info.Mode().String() +
				" 1 ftp ftp " +
				" " + strconv.Itoa(int(info.Size())) + " " +
				info.ModTime().Format(" Jan _2 15:04 ") +
				info.Name() + "\r\n"
		}
	}
	conn.writeMessage(150, "Opening ASCII mode data connection for file list")
	conn.sendByteData([]byte(content))
}

func commandMDTM(conn *ftpConn, param string) {
	path := strings.TrimSpace(param)
	if len(path) < 1 {
		conn.writeMessage(500, "Syntax error")
		return
	}

	fileInfo, err := fsStore.FileInfo(conn.buildPath(path))
	if err != nil {
		conn.writeMessage(550, "File not available")
		return
	} else if !fileInfo.Mode().IsRegular() {
		conn.writeMessage(550, "No such file")
		return
	}

	conn.writeMessage(213, fileInfo.ModTime().Format("20060102150405"))
}

func commandMKD(conn *ftpConn, param string) {
	path := strings.TrimSpace(param)
	if len(path) < 1 {
		conn.writeMessage(500, "Syntax error")
		return
	}
	path = conn.buildPath(path)

	if fileInfo, _ := fsStore.FileInfo(path); fileInfo != nil {
		conn.writeMessage(550, "File or directory already exists")
		return
	}

	if err := fsStore.CreateDir(path); err != nil {
		if os.IsPermission(err) {
			conn.writeMessage(550, "Permission denied")
		} else {
			conn.writeMessage(550, "Directory not available")
		}
		return
	}
	conn.writeMessage(250, "Create directory successfully")
}

func commandMLSD(conn *ftpConn, param string) {
	conn.dataLock.Lock()
	defer conn.dataLock.Unlock()
	if conn.dataConn == nil {
		conn.writeMessage(425, "Error opening data socket")
		return
	}

	path := strings.TrimSpace(param)
	fileInfo, err := fsStore.FileInfo(conn.buildPath(path))
	if err != nil {
		conn.writeMessage(550, "File or directory not available")
		return
	}

	content := ""
	if fileInfo.Mode().IsRegular() {
		content = formatMLSxFileInfo(fileInfo)
	} else {
		infos, err := fsStore.ListDir(conn.buildPath(path))
		if err != nil {
			conn.writeMessage(550, "Directory not available")
			return
		}

		for _, info := range infos {
			content += formatMLSxFileInfo(info)
		}
	}
	conn.writeMessage(150, "Opening ASCII mode data connection for file list")
	conn.sendByteData([]byte(content))
}

func commandMLST(conn *ftpConn, param string) {
	path := strings.TrimSpace(param)
	if len(path) < 1 {
		conn.writeMessage(500, "Syntax error")
		return
	}

	fileInfo, err := fsStore.FileInfo(conn.buildPath(path))
	if err != nil {
		conn.writeMessage(550, "File or directory not available")
		return
	}

	content := "MLST result:\r\n " + formatMLSxFileInfo(fileInfo)
	conn.writeMessageMultiline(250, content)
}

func commandNOOP(conn *ftpConn, param string) {
	conn.writeMessage(200, "OK")
}

func commandOPTS(conn *ftpConn, param string) {
	fields := strings.Fields(param)
	if len(fields) != 2 {
		conn.writeMessage(500, "Syntax error")
		return
	}

	if strings.ToUpper(fields[0]) != "UTF8" {
		conn.writeMessage(501, "Unknown option")
		return
	}

	if strings.ToUpper(fields[1]) == "ON" {
		conn.writeMessage(200, "UTF8 mode enabled")
	} else {
		conn.writeMessage(501, "Unsupported non-utf8 mode")
	}
}

func commandPASS(conn *ftpConn, param string) {
	conn.writeMessage(230, "Guest login ok")
}

func commandPASV(conn *ftpConn, param string) {
	conn.dataLock.Lock()
	if conn.dataConn != nil {
		conn.dataConn.Close()
		conn.dataConn = nil
	}
	conn.dataLock.Unlock()

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
	conn.close()
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
			conn.writeMessage(550, "Permission denied")
		} else {
			conn.writeMessage(550, "File not available")
		}
		return
	}

	conn.writeMessage(150, "Sending file")
	conn.sendStreamData(file)
}

func commandRMD(conn *ftpConn, param string) {
	path := strings.TrimSpace(param)
	if len(path) < 1 {
		conn.writeMessage(500, "Syntax error")
		return
	}
	path = conn.buildPath(path)

	if !fsStore.IsDir(path) {
		conn.writeMessage(550, "No such directory")
		return
	}

	err := fsStore.Delete(path)
	if err != nil {
		conn.writeMessage(550, err.Error())
		return
	}

	conn.writeMessage(250, "Delete directory successfully")
}

func commandSIZE(conn *ftpConn, param string) {
	path := strings.TrimSpace(param)
	if len(path) < 1 {
		conn.writeMessage(500, "Syntax error")
		return
	}

	fileInfo, err := fsStore.FileInfo(conn.buildPath(path))
	if err != nil {
		conn.writeMessage(550, "File not available")
		return
	} else if !fileInfo.Mode().IsRegular() {
		conn.writeMessage(550, "No such file")
		return
	}

	conn.writeMessage(213, strconv.FormatInt(fileInfo.Size(), 10))
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
			conn.writeMessage(550, "Permission denied")
		} else {
			conn.writeMessage(550, "File not available")
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

func formatMLSxFileInfo(file os.FileInfo) string {
	if file.IsDir() {
		return "Type=dir" +
			";Perm=" + file.Mode().String() +
			";Modify=" + file.ModTime().Format("20060102150405") +
			"; " + file.Name() + "\r\n"
	}

	return "Type=file" +
		";Perm=" + file.Mode().String() +
		";Modify=" + file.ModTime().Format("20060102150405") +
		";Size=" + strconv.Itoa(int(file.Size())) +
		"; " + file.Name() + "\r\n"
}
