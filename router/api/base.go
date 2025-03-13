package api

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync/atomic"

	"github.com/whoisnian/glb/tasklane"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/share-Go/global"
)

var (
	_taskLane *tasklane.TaskLane
	_taskSeq  *atomic.Uint64
	_laneMark string
)

func Setup(ctx context.Context) {
	info, err := os.Stat(global.CFG.RootPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(global.CFG.RootPath, osutil.DefaultDirMode)
	}
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		panic("root path is not a directory")
	}
	_taskLane = tasklane.New(ctx, 2, 16)
	_taskSeq = new(atomic.Uint64)
	buf := make([]byte, 5)
	rand.Read(buf)
	_laneMark = base32.StdEncoding.EncodeToString(buf) + "-"
}

func writeNewFile(name string, src io.Reader) error {
	temp := name + ".part"
	file, err := os.OpenFile(temp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, osutil.DefaultFileMode)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, src)
	if err1 := file.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err2 := os.Rename(temp, name); err2 != nil && err == nil {
		err = err2
	}
	return err
}

type downloadTask struct {
	tid string
	url string
	dir string
}

func newDownloadTask(url string, dir string) *downloadTask {
	return &downloadTask{
		tid: _laneMark + strconv.FormatUint(_taskSeq.Add(1), 10),
		url: url,
		dir: dir,
	}
}

func (t *downloadTask) Start() {
	global.LOG.Infof(context.TODO(), "task(%s) start download %s", t.tid, t.url)
	resp, err := http.Get(t.url)
	if err != nil {
		global.LOG.Errorf(context.TODO(), "task(%s) download failed: %v", t.tid, err)
		return
	}
	defer resp.Body.Close()
	fpath := filepath.Join(t.dir, path.Base(t.url))
	if err := writeNewFile(fpath, resp.Body); err != nil {
		global.LOG.Errorf(context.TODO(), "task(%s) writeNewFile failed: %v", t.tid, err)
		return
	}
	global.LOG.Infof(context.TODO(), "task(%s) save as %s done", t.tid, fpath)
}

type respMessage struct {
	Message string
}

const (
	typeRegular   int64 = 0
	typeDirectory int64 = 1
)

type respFileInfo struct {
	Type  int64
	Name  string
	Size  int64
	MTime int64
}

type respFileInfos struct {
	FileInfos []respFileInfo
}

func prepareRespFileInfo(info os.FileInfo) respFileInfo {
	t := typeRegular
	if info.Mode().IsDir() {
		t = typeDirectory
	}
	return respFileInfo{
		Type:  t,
		Name:  info.Name(),
		Size:  info.Size(),
		MTime: info.ModTime().Unix(),
	}
}
