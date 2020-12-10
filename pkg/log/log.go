package log

import (
	logger "log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/whoisnian/share-Go/pkg/state"
)

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// MakeHander ...
func MakeHander(fn func(store state.Store)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		store := state.NewStore(w, r)
		fn(store)
		logger.Printf("%s [%d] %s %s %s %d",
			r.RemoteAddr[0:strings.IndexByte(r.RemoteAddr, ':')],
			store.Code,
			r.Method,
			r.URL.Path,
			r.UserAgent(),
			time.Now().Sub(start).Milliseconds())
	}
}
