package httpd

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/whoisnian/share-Go/pkg/logger"
	"github.com/whoisnian/share-Go/pkg/util"
)

const routeAny string = "/*"
const routeParam string = "/:param"

var methodList = map[string]string{
	"GET":     "/get",
	"HEAD":    "/head",
	"POST":    "/post",
	"PUT":     "/put",
	"DELETE":  "/delete",
	"CONNECT": "/connect",
	"OPTIONS": "/options",
	"TRACE":   "/trace",
	"PATCH":   "/patch",
}

type nodeData struct {
	route         string
	handler       func(Store)
	paramNameList []string
}

type routeNode struct {
	next map[string]*routeNode
	data *nodeData
}

func (node *routeNode) nextNode(name string) (res *routeNode) {
	if res, ok := node.next[name]; ok {
		return res
	}
	if node.next == nil {
		node.next = make(map[string]*routeNode)
	}
	res = new(routeNode)
	node.next[name] = res
	return res
}

func parseRoute(node *routeNode, route string) (*routeNode, []string) {
	var paramNameList []string
	fragments := strings.Split(route, "/")
	for _, fragment := range fragments {
		if len(fragment) < 1 {
			continue
		} else if fragment == "*" {
			paramNameList = append(paramNameList, routeAny)
			node = node.nextNode(routeAny)
		} else if fragment[0] == ':' {
			paramName := fragment[1:]
			if paramName == "" || util.Contain(paramNameList, paramName) {
				logger.Fatal("Invalid fragment '", fragment, "' for route: '", route, "'")
			}
			paramNameList = append(paramNameList, paramName)
			node = node.nextNode(routeParam)
		} else {
			node = node.nextNode(fragment)
		}
	}
	return node, paramNameList
}

func findRoute(node *routeNode, route string) (*routeNode, []string) {
	var paramValueList []string
	fragments := strings.Split(route, "/")
	for index, fragment := range fragments {
		if len(fragment) < 1 {
			continue
		} else if res, ok := node.next[fragment]; ok {
			node = res
		} else if res, ok := node.next[routeParam]; ok {
			paramValueList = append(paramValueList, fragment)
			node = res
		} else if res, ok := node.next[routeAny]; ok {
			paramValueList = append(paramValueList, strings.Join(fragments[index:], "/"))
			node = res
			break
		} else {
			return nil, nil
		}
	}
	return node, paramValueList
}

type serveMux struct {
	mu   sync.RWMutex
	root *routeNode
}

func (mux *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	store := Store{&statusResponseWriter{w, http.StatusOK}, r, make(map[string]string)}

	defer func() {
		if err := recover(); err != nil {
			store.Error500("Internal Server Error")
		}

		logger.Info(
			r.RemoteAddr[0:strings.IndexByte(r.RemoteAddr, ':')], " [",
			store.w.status, "] ",
			r.Method, " ",
			r.URL.Path, " ",
			r.UserAgent(), " ",
			time.Now().Sub(start).Milliseconds(),
		)
	}()

	methodTag, ok := methodList[r.Method]
	if !ok {
		store.Respond404()
		return
	}
	node, paramValueList := findRoute(mux.root, r.URL.EscapedPath())
	if node == nil {
		store.Respond404()
		return
	}

	res, ok := node.next[methodTag]
	if !ok {
		store.Respond404()
		return
	}

	for index, paramName := range res.data.paramNameList {
		store.m[paramName] = paramValueList[index]
	}

	res.data.handler(store)
}

var mux *serveMux

func init() {
	mux = new(serveMux)
	mux.root = new(routeNode)
}

// Handle registers the handler for the given route.
func Handle(route string, method string, handler func(Store)) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	methodTag, ok := methodList[method]
	if !ok {
		logger.Fatal("Invalid method '", method, "' for route: '", route, "'")
	}

	node, paramNameList := parseRoute(mux.root, route)
	if _, ok = node.next[methodTag]; ok {
		logger.Fatal("Duplicate method '", method, "' for route: '", route, "'")
	}
	node.nextNode(methodTag).data = &nodeData{route, handler, paramNameList}
}

// Start listens on the addr and then creates goroutine to handle each request.
func Start(addr string) {
	logger.Info("Service httpd started: <http://", addr, ">")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Fatal(err)
	}
}
