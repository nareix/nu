package main

import (
	"net/http"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/nareix/nu/pack"
	"github.com/nareix/nu/utils"
)

type HttpCommand struct {
	Secure bool   `arg:"-s"`
	Addr   string `arg:"-a"`
	Cert   string `arg:"-c"`
	Key    string `arg:"-k"`
	Dir    string `arg:"-d"`
}

func NewHttpCommand() *HttpCommand {
	return &HttpCommand{
		Addr: ":8080",
		Dir:  ".",
	}
}

func (c *HttpCommand) Run() error {
	if c.Secure {
		return http.ListenAndServeTLS(
			c.Addr, c.Cert, c.Key,
			http.FileServer(http.Dir(c.Dir)),
		)
	} else {
		return http.ListenAndServe(c.Addr, http.FileServer(http.Dir(c.Dir)))
	}
}

type RootCmmand struct {
	Http *HttpCommand  `arg:"subcommand:http"`
	Pack *pack.Command `arg:"subcommand:pack"`
}

func (c *RootCmmand) Run() error {
	switch {
	case c.Http != nil:
		return c.Http.Run()
	case c.Pack != nil:
		return c.Pack.Run()
	default:
		return arg.ErrHelp
	}
}

func mainfn() error {
	root := &RootCmmand{
		Http: NewHttpCommand(),
	}
	p := arg.MustParse(root)
	if err := root.Run(); err != nil {
		if err == arg.ErrHelp {
			p.WriteHelp(os.Stderr)
			return nil
		}
		return err
	}
	return nil
}

func main() {
	utils.MustRun(mainfn)
}
