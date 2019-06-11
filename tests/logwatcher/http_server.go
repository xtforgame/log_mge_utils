// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logwatcher

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	// funk "github.com/thoas/go-funk"
	"github.com/xtforgame/log_mge_utils/fshelper"
	"github.com/xtforgame/log_mge_utils/httpserver"
	"net/http"
	// "sort"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	fmt.Println("path :", path)
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

type WatcherStats struct {
	Name    string `json:"name,"`
	Logs    int64  `json:"logs,omitempty"`
	Error   string `json:"error,omitempty"`
	LogSize int64  `json:"logSize,omitempty"`
}

type HttpServer struct {
	logPath string
	webPath string
	server  *http.Server
	router  *chi.Mux
}

func NewHttpServer(logPath string, webPath string) *HttpServer {
	r := chi.NewRouter()
	return &HttpServer{
		logPath: logPath,
		webPath: webPath,
		server: &http.Server{
			Addr:    ":8080",
			Handler: r,
		},
		router: r,
	}
}

func GetWatcherInfo(watcherPath string, watcherName string) (*WatcherStats, error) {
	logList, err := fshelper.ListDir(watcherPath)
	if err != nil {
		return nil, err
	}
	var totalSize int64
	for _, fileInfo := range logList.Files {
		totalSize += fileInfo.Size()
	}
	return &WatcherStats{
		Name:    watcherName,
		Logs:    int64(len(logList.Files)),
		LogSize: totalSize,
	}, nil
}

func (hs *HttpServer) Init() {
	if LoggerHeplerInst == nil {
		LoggerHeplerInst = CreateLoggerHepler(hs.logPath)
	}
	hs.router.HandleFunc("/logger", hs.loggerHome)
	hs.router.HandleFunc("/listener", hs.listenerHome)
	hs.router.HandleFunc("/logger/{logID}", hs.loggerHome)
	hs.router.HandleFunc("/listener/{logID}", hs.listenerHome)
	hs.router.HandleFunc("/app.js", hs.jsScript)
	hs.router.HandleFunc("/loggers/{logID}", LoggerWebsocket)
	hs.router.HandleFunc("/listeners/{logID}", ListenerWebsocket)
	// hs.router.FileServer("/", http.Dir("web/"))
	// FileServer(hs.router, "/assets", http.Dir("./assets"))
	hs.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	hs.router.Route("/stats", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			watchersFolder := filepath.Join(hs.logPath, "log-watcher")
			watcherList, err := fshelper.ListDir(watchersFolder)
			stats := []*WatcherStats{}
			for _, fileInfo := range watcherList.Dirs {
				watcherStats, err := GetWatcherInfo(filepath.Join(watchersFolder, fileInfo.Name()), fileInfo.Name())
				if err != nil {
					stats = append(stats, &WatcherStats{
						Name:  fileInfo.Name(),
						Error: err.Error(),
					})
				} else {
					stats = append(stats, watcherStats)
				}
			}
			if err != nil {
				fmt.Println("err :", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Log directory is not accessible"))
				return
			}
			if jsonBytes, err := json.Marshal(stats); err == nil {
				w.Write(jsonBytes)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Server Error"))
				return
			}
			// w.Write([]byte("[]"))
		})

		r.Get("/{logID:[0-9a-zA-Z_-]+}", func(w http.ResponseWriter, r *http.Request) {
			logID := chi.URLParam(r, "logID")
			watcherFolder := filepath.Join(hs.logPath, "log-watcher", logID)
			watcherStats, err := GetWatcherInfo(watcherFolder, logID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Not Found"))
			} else {
				if jsonBytes, err := json.Marshal(watcherStats); err == nil {
					w.Write(jsonBytes)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - Internal Server Error"))
					return
				}
				return
			}
		})

		r.Delete("/{logID:[0-9a-zA-Z_-]+}", func(w http.ResponseWriter, r *http.Request) {
			logID := chi.URLParam(r, "logID")
			LoggerHeplerInst.ForceRemoveLogger(logID)
			w.Write([]byte("{\"result\":\"done\"}"))
		})
	})
}

func (hs *HttpServer) Start() {
	loggerHomeHtmlTmp, _ := ioutil.ReadFile(filepath.Join(hs.webPath, "logwatcher/logger.html"))
	loggerHomeTemplate = template.Must(template.New("").Parse(string(loggerHomeHtmlTmp)))

	listenerHomeHtmlTmp, _ := ioutil.ReadFile(filepath.Join(hs.webPath, "logwatcher/listener.html"))
	listenerHomeTemplate = template.Must(template.New("").Parse(string(listenerHomeHtmlTmp)))

	jsTmp, _ := ioutil.ReadFile(filepath.Join(hs.webPath, "logwatcher/app.js"))
	jsTemplate = template.Must(template.New("").Parse(string(jsTmp)))
	httpserver.RunAndWaitGracefulShutdown(hs.server)
}

func (hs *HttpServer) loggerHome(w http.ResponseWriter, r *http.Request) {
	logID := chi.URLParam(r, "logID")
	if logID == "" {
		logID = "20022"
	}
	// fmt.Println("logID :", logID)
	/* ======================= for test start ======================= */
	loggerHomeHtmlTmp, _ := ioutil.ReadFile(filepath.Join(hs.webPath, "logwatcher/logger.html"))
	loggerHomeTemplate = template.Must(template.New("").Parse(string(loggerHomeHtmlTmp)))
	/* =======================  for test end  ======================= */
	loggerHomeTemplate.Execute(
		w,
		struct {
			WsBaseUrl              string
			LogID                  string
			EventNextIterationCode byte
			EventOnDataCode        byte
			EventLogRemovedCode    byte
		}{
			"ws://" + r.Host,
			logID,
			EventNextIterationCode,
			EventOnDataCode,
			EventLogRemovedCode,
		},
	)
}

func (hs *HttpServer) listenerHome(w http.ResponseWriter, r *http.Request) {
	logID := chi.URLParam(r, "logID")
	if logID == "" {
		logID = "20022"
	}
	// fmt.Println("logID :", logID)
	/* ======================= for test start ======================= */
	listenerHomeHtmlTmp, _ := ioutil.ReadFile(filepath.Join(hs.webPath, "logwatcher/listener.html"))
	listenerHomeTemplate = template.Must(template.New("").Parse(string(listenerHomeHtmlTmp)))
	/* =======================  for test end  ======================= */
	listenerHomeTemplate.Execute(
		w,
		struct {
			WsBaseUrl              string
			LogID                  string
			EventNextIterationCode byte
			EventOnDataCode        byte
			EventLogRemovedCode    byte
		}{
			"ws://" + r.Host,
			logID,
			EventNextIterationCode,
			EventOnDataCode,
			EventLogRemovedCode,
		},
	)
}

func (hs *HttpServer) jsScript(w http.ResponseWriter, r *http.Request) {
	/* ======================= for test start ======================= */
	jsTmp, _ := ioutil.ReadFile(filepath.Join(hs.webPath, "logwatcher/app.js"))
	jsTemplate = template.Must(template.New("").Parse(string(jsTmp)))
	/* =======================  for test end  ======================= */
	jsTemplate.Execute(
		w,
		struct {
			WsBaseUrl              string
			EventNextIterationCode byte
			EventOnDataCode        byte
			EventLogRemovedCode    byte
		}{
			"ws://" + r.Host,
			EventNextIterationCode,
			EventOnDataCode,
			EventLogRemovedCode,
		},
	)
}

var loggerHomeTemplate *template.Template
var listenerHomeTemplate *template.Template
var jsTemplate *template.Template
