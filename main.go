package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/nareix/nu/fiximport"
	"github.com/nareix/nu/pack"
	"github.com/nareix/nu/utils"
	"github.com/nareix/nu/utils/go-arg"
)

type HttpCommand struct {
	Secure            bool   `arg:"-s"`
	Addr              string `arg:"-a"`
	Cert              string `arg:"-c"`
	Key               string `arg:"-k"`
	Dir               string `arg:"-d"`
	HostJSFile        string `arg:"--js"`
	HostOnlyIndexHTML bool   `arg:"--indexhtml"`
}

func (c *HttpCommand) Run() error {
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.Dir == "" {
		c.Dir = "."
	}
	log.Println("server http on", c.Addr, "dir", c.Dir)

	serve := func(handler http.Handler) error {
		if c.Secure {
			return http.ListenAndServeTLS(
				c.Addr, c.Cert, c.Key,
				handler,
			)
		} else {
			return http.ListenAndServe(c.Addr, handler)
		}
	}

	if c.HostJSFile != "" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<html><head><script src="index.js"></script></head><body></body></html>`))
		})
		http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, c.HostJSFile)
		})
		return serve(nil)
	} else if c.HostOnlyIndexHTML {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/static") {
				http.FileServer(http.Dir(c.Dir)).ServeHTTP(w, r)
				return
			}
			http.ServeFile(w, r, filepath.Join(c.Dir, "index.html"))
		})
		return serve(nil)
	} else {
		return serve(http.FileServer(http.Dir(c.Dir)))
	}
}

type RootCmmand struct {
	Http    *HttpCommand       `arg:"subcommand:http"`
	Pack    *pack.Command      `arg:"subcommand:pack"`
	Rewrite *fiximport.Command `arg:"subcommand:rewrite"`
}

func mainfn() error {
	root := &RootCmmand{}
	return arg.MustRun(root)
}

func main() {
	utils.MustRun(mainfn)
}
