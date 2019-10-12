package pack

import (
	"archive/tar"
	"compress/gzip"
	"debug/elf"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/flynn/json5"
)

// func Exists(name string) bool {
// 	if _, err := os.Stat(name); err != nil {
// 		if os.IsNotExist(err) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func Realpath(fpath string) (string, error) {
// 	if len(fpath) == 0 {
// 		return "", os.ErrInvalid
// 	}

// 	if !filepath.IsAbs(fpath) {
// 		pwd, err := os.Getwd()
// 		if err != nil {
// 			return "", err
// 		}
// 		fpath = filepath.Join(pwd, fpath)
// 	}

// 	path := []byte(fpath)
// 	nlinks := 0
// 	start := 1
// 	prev := 1
// 	for start < len(path) {
// 		c := nextComponent(path, start)
// 		cur := c[start:]

// 		switch {

// 		case len(cur) == 0:
// 			copy(path[start:], path[start+1:])
// 			path = path[0 : len(path)-1]

// 		case len(cur) == 1 && cur[0] == '.':
// 			if start+2 < len(path) {
// 				copy(path[start:], path[start+2:])
// 			}
// 			path = path[0 : len(path)-2]

// 		case len(cur) == 2 && cur[0] == '.' && cur[1] == '.':
// 			copy(path[prev:], path[start+2:])
// 			path = path[0 : len(path)+prev-(start+2)]
// 			prev = 1
// 			start = 1

// 		default:

// 			fi, err := os.Lstat(string(c))
// 			if err != nil {
// 				return "", err
// 			}
// 			if isSymlink(fi) {

// 				nlinks++
// 				if nlinks > 16 {
// 					return "", os.ErrInvalid
// 				}

// 				var link string
// 				link, err = os.Readlink(string(c))
// 				after := string(path[len(c):])

// 				// switch symlink component with its real path
// 				path = switchSymlinkCom(path, start, link, after)

// 				prev = 1
// 				start = 1
// 			} else {
// 				// Directories
// 				prev = start
// 				start = len(c) + 1
// 			}
// 		}
// 	}

// 	for len(path) > 1 && path[len(path)-1] == os.PathSeparator {
// 		path = path[0 : len(path)-1]
// 	}
// 	return string(path), nil
// }

// // test if a link is symbolic link
// func isSymlink(fi os.FileInfo) bool {
// 	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
// }

// // switch a symbolic link component to its real path
// func switchSymlinkCom(path []byte, start int, link, after string) []byte {

// 	if link[0] == os.PathSeparator {
// 		// Absolute links
// 		return []byte(filepath.Join(link, after))
// 	}

// 	// Relative links
// 	return []byte(filepath.Join(string(path[0:start]), link, after))
// }

// // get the next component
// func nextComponent(path []byte, start int) []byte {
// 	v := bytes.IndexByte(path[start:], os.PathSeparator)
// 	if v < 0 {
// 		return path
// 	}
// 	return path[0 : start+v]
// }

// var ldlinuxpath = "/lib64/ld-linux-x86-64.so.2"

// func otoolL(lib string) (paths []string, err error) {
// 	c := exec.Command("otool", "-L", lib)
// 	stdout, _ := c.StdoutPipe()
// 	br := bufio.NewReader(stdout)
// 	if err = c.Start(); err != nil {
// 		err = fmt.Errorf("otoolL: %s", err)
// 		return
// 	}
// 	for i := 0; ; i++ {
// 		var line string
// 		var rerr error
// 		if line, rerr = br.ReadString('\n'); rerr != nil {
// 			break
// 		}
// 		if i == 0 {
// 			continue
// 		}
// 		f := strings.Fields(line)
// 		if len(f) >= 2 && strings.HasPrefix(f[1], "(") {
// 			paths = append(paths, f[0])
// 		}
// 	}
// 	return
// }

// func installNameTool(lib string, change [][]string) error {
// 	if len(change) == 0 {
// 		return nil
// 	}
// 	args := []string{}
// 	for _, c := range change {
// 		args = append(args, "-change")
// 		args = append(args, c[0])
// 		args = append(args, c[1])
// 	}
// 	args = append(args, lib)

// 	return runcmd("install_name_tool", args...)
// }

// type CopyEntry struct {
// 	Realpath string
// 	IsBin    bool
// }

// func packlibDarwin(copy []CopyEntry) error {
// 	visited := map[string]bool{}
// 	isbin := map[string]bool{}

// 	var dfs func(k string) error
// 	dfs = func(k string) error {
// 		if visited[k] {
// 			return nil
// 		}
// 		visited[k] = true
// 		paths, err := otoolL(k)
// 		if err != nil {
// 			return err
// 		}
// 		for _, p := range paths {
// 			if strings.HasPrefix(p, "@") {
// 				const rpath = "@rpath/"
// 				if strings.HasPrefix(p, rpath) {
// 					lp := locate(strings.TrimPrefix(p, rpath))
// 					if lp != "" {
// 						p = lp
// 					} else {
// 						continue
// 					}
// 				} else {
// 					continue
// 				}
// 			}
// 			if strings.HasPrefix(p, "/usr/lib") {
// 				continue
// 			}
// 			if strings.HasPrefix(p, "/System") {
// 				continue
// 			}
// 			if err := dfs(p); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}

// 	for _, f := range copy {
// 		if f.IsBin {
// 			isbin[f.Realpath] = true
// 		}
// 		if err := dfs(f.Realpath); err != nil {
// 			return err
// 		}
// 	}

// 	change := [][]string{}
// 	for p := range visited {
// 		if !strings.HasPrefix(p, "/") {
// 			continue
// 		}
// 		fname := path.Join("lib", path.Base(p))
// 		change = append(change, []string{p, fname})
// 	}

// 	for p := range visited {
// 		dstdir := "lib"
// 		if isbin[p] {
// 			dstdir = "bin"
// 		}
// 		fname := path.Join(dstdir, path.Base(p))
// 		if err := runcmd("cp", "-f", p, fname); err != nil {
// 			return err
// 		}
// 		if err := runcmd("chmod", "744", fname); err != nil {
// 			return err
// 		}
// 		if err := installNameTool(fname, change); err != nil {
// 			return err
// 		}
// 		fmt.Println("copy", fname)
// 	}

// 	return nil
// }

// var libsearchpath []string

// func locate(name string) (out string) {
// 	for _, root := range libsearchpath {
// 		p := path.Join(root, name)
// 		_, serr := os.Stat(p)
// 		if serr == nil {
// 			out = p
// 			return
// 		}
// 	}
// 	return
// }

// type Entry struct {
// 	Name     string
// 	Realpath string
// 	IsBin    bool
// }

type Config struct {
	LibSearchPath []string `json:"libSearchPath"`
	Lib           []string `json:"lib"`
	Bin           []string `json:"bin"`
	Etc           []string `json:"etc"`
}

func ELFFileImportLibs(filename string) ([]string, error) {
	f, err := elf.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.ImportedLibraries()
}

func SearchLib(searchpath []string, filename string) (string, bool) {
	for _, dir := range searchpath {
		searchfilename := path.Join(dir, filename)
		if s, err := os.Stat(searchfilename); err == nil && !s.IsDir() {
			return searchfilename, true
		}
	}
	return "", false
}

type LibWalker struct {
	visited map[string]struct{}
}

func NewLibWalker() *LibWalker {
	return &LibWalker{
		visited: map[string]struct{}{},
	}
}

func (w *LibWalker) Walk(searchpath []string, filename string, fn func(filename, realpath string) error) error {
	filerealpath, ok := SearchLib(searchpath, filename)
	if !ok {
		return fmt.Errorf("lib %s not found in search paths %s", filename, searchpath)
	}
	if _, ok := w.visited[filename]; ok {
		return nil
	}
	w.visited[filename] = struct{}{}
	if err := fn(filename, filerealpath); err != nil {
		return err
	}
	importlibs, err := ELFFileImportLibs(filerealpath)
	if err != nil {
		return err
	}
	for _, filename := range importlibs {
		if err := w.Walk(searchpath, filename, fn); err != nil {
			return err
		}
	}
	return nil
}

type TarSaver struct {
	f  *os.File
	gw *gzip.Writer
	tw *tar.Writer
}

func (s *TarSaver) Close() error {
	if err := s.tw.Close(); err != nil {
		return err
	}
	if err := s.gw.Close(); err != nil {
		return err
	}
	if err := s.f.Close(); err != nil {
		return err
	}
	return nil
}

func (s *TarSaver) WriteFile(tarpath, filename string) error {
	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(stat, "")
	if err != nil {
		return err
	}
	header.Name = tarpath
	if err := s.tw.WriteHeader(header); err != nil {
		return err
	}
	if !stat.IsDir() {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(s.tw, file); err != nil {
			return err
		}
	}
	return nil
}

func CreateTarSaver(filename string) (*TarSaver, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	s := &TarSaver{
		gw: gw,
		tw: tw,
		f:  f,
	}
	return s, nil
}

func (c *Command) Pack(config Config) error {
	var savetotar *TarSaver
	if c.SaveTo != "" {
		ts, err := CreateTarSaver(c.SaveTo)
		if err != nil {
			return err
		}
		defer ts.Close()
		savetotar = ts
	}

	walksavefile := func(savepath, filename string) error {
		if savetotar != nil {
			if err := savetotar.WriteFile(savepath, filename); err != nil {
				return err
			}
		} else {
			fmt.Println(savepath)
		}
		return nil
	}

	walker := NewLibWalker()
	walklib := func(filename, realpath string) error {
		savepath := path.Join("lib", filename)
		return walksavefile(savepath, realpath)
	}

	for _, filename := range config.Lib {
		filenames, err := filepath.Glob(filename)
		if err != nil {
			return err
		}
		for _, filename := range filenames {
			filename = path.Base(filename)
			if err := walker.Walk(config.LibSearchPath, filename, walklib); err != nil {
				return err
			}
		}
	}

	for _, filename := range config.Bin {
		filenames, err := filepath.Glob(filename)
		if err != nil {
			return err
		}
		for _, filename := range filenames {
			savepath := path.Join("bin", filename)
			if err := walksavefile(savepath, filename); err != nil {
				return err
			}
			importlibs, err := ELFFileImportLibs(filename)
			if err != nil {
				return err
			}
			for _, filename := range importlibs {
				if err := walker.Walk(config.LibSearchPath, filename, walklib); err != nil {
					return err
				}
			}
		}
	}

	for _, filename := range config.Etc {
		filenames, err := filepath.Glob(filename)
		if err != nil {
			return err
		}
		for _, filename := range filenames {
			filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
				if err := walksavefile(path, path); err != nil {
					return err
				}
				return nil
			})
		}
	}

	return nil
}

type Command struct {
	ConfigFile string `arg:"positional"`
	SaveTo     string `arg:"-w"`
}

func (c *Command) Run() error {
	configfilebytes, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		return err
	}
	config := Config{}
	if err := json5.Unmarshal(configfilebytes, &config); err != nil {
		return err
	}
	return c.Pack(config)
}
