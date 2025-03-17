package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/whoisnian/glb/tasklane"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/share-Go/global"
)

var (
	_taskLane *tasklane.TaskLane
	_taskSeq  *tasklane.Sequence
)

func Setup(ctx context.Context) {
	info, err := os.Stat(global.CFG.RootPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(global.CFG.RootPath, osutil.DefaultDirMode)
		}
	} else if !info.IsDir() {
		err = errors.New("root path is not a directory")
	}
	if err != nil {
		panic(err)
	}
	_taskLane = tasklane.New(ctx, 2, 16)
	_taskSeq = tasklane.NewSequence(8)
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
		tid: string(_taskSeq.Next()),
		url: url,
		dir: dir,
	}
}

func (t *downloadTask) Start(ctx context.Context) {
	global.LOG.Infof(ctx, "task(%s) start download %s", t.tid, t.url)
	resp, err := http.Get(t.url)
	if err != nil {
		global.LOG.Errorf(ctx, "task(%s) download failed: %v", t.tid, err)
		return
	}
	defer resp.Body.Close()
	fpath := filepath.Join(t.dir, path.Base(t.url))
	if err := writeNewFile(fpath, resp.Body); err != nil {
		global.LOG.Errorf(ctx, "task(%s) writeNewFile failed: %v", t.tid, err)
		return
	}
	global.LOG.Infof(ctx, "task(%s) save as %s done", t.tid, fpath)
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
