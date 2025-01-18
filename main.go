package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/webx-top/com"
	"golang.org/x/mod/modfile"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`usage:`, os.Args[0], `go.mod`)
		return
	}

	modfilePath := os.Args[1] // `/Users/hank/go/src/github.com/admpub/nging/go.mod`
	data, err := os.ReadFile(modfilePath)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := modfile.Parse(`in`, data, nil)
	if err != nil {
		log.Println(err)
		return
	}
	pkgDir := filepath.Dir(modfilePath)
	vendorDir := filepath.Join(pkgDir, `vendor`)
	results := Files{}
	for _, req := range f.Require {
		dir := filepath.Join(vendorDir, req.Mod.Path)
		fi, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			log.Println(err)
			return
		}
		if !fi.IsDir() {
			continue
		}
		size, err := dirSize(dir)
		if err != nil {
			log.Println(err)
			return
		}
		results = append(results, &File{Path: req.Mod.Path, Size: size})
	}
	sort.Sort(results)
	for _, file := range results {
		fmt.Println(file.Path + "\t\t\t" + com.HumaneFileSize(uint64(file.Size)))
	}
}

func dirSize(dir string) (int64, error) {
	var size int64
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += info.Size()
		return err
	})
	return size, err
}

type File struct {
	Path string
	Size int64
}

type Files []*File

func (f Files) Len() int {
	return len(f)
}

func (f Files) Less(i, j int) bool {
	return f[i].Size < f[j].Size
}
func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
