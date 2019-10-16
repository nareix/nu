package fiximport

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Command struct {
	Paths          []string `arg:"positional"`
	RewriteImports []string `arg:"-i,separate"`
	RewritePackage string   `arg:"-p"`
}

func (c *Command) Run() error {
	if len(c.Paths) == 0 {
		c.Paths = []string{"."}
	}

	log.Println("rewrite", c.RewriteImports)

	if len(c.RewriteImports)%2 != 0 {
		return fmt.Errorf("rewrite imports must be pair")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	walkfn := func(walkpath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		src, err := ioutil.ReadFile(walkpath)
		if err != nil {
			return err
		}

		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
		if err != nil {
			return nil
		}

		if s := c.RewritePackage; s != "" {
			if s == "." {
				s = path.Base(cwd)
			}
			f.Name = ast.NewIdent(s)
		}

		for _, imp := range f.Imports {
			imppath := strings.Trim(imp.Path.Value, `"`)
			replace := func(match, rep string) {
				if strings.HasPrefix(imppath, match) {
					newimppath := rep + imppath[len(match):]
					imp.Path.Value = `"` + newimppath + `"`
				}
			}

			for i := 0; i < len(c.RewriteImports); i += 2 {
				before := c.RewriteImports[i]
				after := c.RewriteImports[i+1]
				replace(before, after)
			}
		}

		fout, err := os.Create(walkpath)
		if err != nil {
			return err
		}
		defer fout.Close()

		if err := format.Node(fout, fset, f); err != nil {
			return err
		}
		return nil
	}

	for _, p := range c.Paths {
		if err := filepath.Walk(p, walkfn); err != nil {
			return err
		}
	}
	return nil
}
