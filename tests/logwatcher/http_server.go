// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logwatcher

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	// funk "github.com/thoas/go-funk"
	"github.com/xtforgame/log_mge_utils/httpserver"
	"net/http"
	// "sort"
	"io/ioutil"
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

type HttpServer struct {
	server *http.Server
	router *chi.Mux
}

func NewHttpServer() *HttpServer {
	r := chi.NewRouter()
	return &HttpServer{
		server: &http.Server{
			Addr:    ":8080",
			Handler: r,
		},
		router: r,
	}
}

func (hs *HttpServer) Init() {
	hs.router.HandleFunc("/client", home)
	hs.router.HandleFunc("/echo", TestHandleWebsocket)
	// hs.router.FileServer("/", http.Dir("web/"))
	// FileServer(hs.router, "/assets", http.Dir("./assets"))
	hs.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	hs.router.Get("/reg", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
}

func (hs *HttpServer) Start() {
	dat, _ := ioutil.ReadFile("./web/logwatcher/client.html")
	homeTemplate = template.Must(template.New("").Parse(string(dat)))
	httpserver.RunAndWaitGracefulShutdown(hs.server)
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate *template.Template
