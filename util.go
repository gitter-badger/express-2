package express

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"
)

var replacer = strings.NewReplacer(":", "", "?", "")

func haveParameters(path string) bool {
	for i := 0; i < len(path); i++ {
		if path[i] == ':' {
			return true
		}
		if path[i] == '*' {
			return true
		}
		if path[i] == '?' {
			return true
		}
	}
	return false
}

func stripParameters(path string) (params []string) {
	segments := strings.Split(path, "/")
	for _, s := range segments {
		if s == "" {
			continue
		}
		if strings.Contains(s, ":") {
			s = replacer.Replace(s)
			params = append(params, s)
			continue
		}
		if strings.Contains(s, "*") {
			params = append(params, "*")
		}
	}
	return params
}

func pathToRegex(path string) (regex string) {
	regex = "^"
	segments := strings.Split(path, "/")
	for _, s := range segments {
		if s == "" {
			continue
		}
		if strings.Contains(s, ":") {
			if strings.Contains(s, "?") {
				regex += "(?:/([^/]+?))?"
			} else {
				regex += "/(?:([^/]+?))"
			}
		} else if strings.Contains(s, "*") {
			regex += "/(.*)"
		} else {
			regex += "/" + s
		}
	}
	regex += "/?$"
	return regex
}

func dirWalk(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// Credits to @savsgio
// https://github.com/savsgio/gotils/blob/master/conv.go

// b2s converts byte slice to a string without memory allocation.
func b2s(b []byte) string {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&b))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*string)(unsafe.Pointer(&bh))
}

// s2b converts string to a byte slice without memory allocation.
func s2b(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
