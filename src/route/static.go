// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a GPLv3 license that can be found in the LICENSE file.

package route

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.design/x/ssaplayground/src/config"
)

// serveFS is a middleware that allows static files serves
// in the root router like "addr:port/""
type serveFS interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func loacalFile(root string, indexes bool) serveFS {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if !l.indexes && stats.IsDir() {
			return false
		}
		return true
	}
	return false
}

// static returns a middleware handler that serves static files
// in the given directory.
func static(urlPrefix string) gin.HandlerFunc {
	fs := loacalFile(config.Get().Static, true)
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
