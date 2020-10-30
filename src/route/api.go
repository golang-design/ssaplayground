// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a GPLv3 license that can be found in the LICENSE file.

package route

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.design/x/ssaplayground/src/config"
)

// PingInput is a a reserved structure
type PingInput struct {
}

// PingOutput is used for service health
type PingOutput struct {
	Message   string `json:"message"`
	GoVersion string `json:"go_version"`
}

// Pong response for health check
func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, &PingOutput{
		Message:   "pong",
		GoVersion: runtime.Version(),
	})
}

// BuildSSAInput ...
type BuildSSAInput struct {
	FuncName string `json:"funcname"`
	GcFlags  string `json:"gcflags"`
	Code     string `json:"code"`
}

// BuildSSAOutput ...
type BuildSSAOutput struct {
	BuildID string `json:"build_id"`
	Msg     string `json:"msg"`
}

// BuildSSA serves the code send by user and builds its SSA IR into html.
// TODO: speedup for request response, e.g. as async rest api.
func BuildSSA(c *gin.Context) {
	// 1. create a folder in config.Get().Static/buildbox
	out := BuildSSAOutput{
		BuildID: uuid.New().String(),
	}
	path := filepath.Join(config.Get().Static, "/buildbox", "/"+out.BuildID)

	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		out.Msg = fmt.Sprintf("cannot create buildbox, err: %v", err)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	// 2. write code
	var in BuildSSAInput
	err = c.BindJSON(&in)
	if err != nil {
		os.Remove(path)
		out.Msg = fmt.Sprintf("cannot bind input params, err: \n%v", err)
		c.JSON(http.StatusInternalServerError, out)
		return
	}
	if !findSSAFunc(in.Code, in.FuncName) {
		os.Remove(path)
		out.Msg = fmt.Sprintf("cannot find GOSSAFUNC=%s in your code.", in.FuncName)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	var buildFile string
	isTest := isPackageTest(in.Code)
	if !isTest {
		buildFile = filepath.Join(path, "/main.go")
	} else {
		buildFile = filepath.Join(path, "/main_test.go")
	}
	err = ioutil.WriteFile(buildFile, []byte(in.Code), os.ModePerm)
	if err != nil {
		os.Remove(path)
		out.Msg = fmt.Sprintf("cannot save your code, err: \n%v", err)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	// 3.1 goimports
	err = autoimports(buildFile)
	if err != nil {
		os.Remove(path)
		out.Msg = fmt.Sprintf("cannot run autoimports for your code, err: \n%v", err)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// 3.2 go mod init gossa && go mod tidy
	err = initModules(path)
	if err != nil {
		os.Remove(path)
		out.Msg = fmt.Sprintf("cannot use go modules for your code, err: \n%v", err)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// 3.3 GOSSAFUNC=foo go build
	outFile := filepath.Join(path, "/main.out")
	err = buildSSA(in.FuncName, in.GcFlags, outFile, buildFile, isTest)
	if err != nil {
		os.Remove(path)
		out.Msg = fmt.Sprintf("cannot build ssa for your code, err: \n%v", err)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	os.Remove(outFile) // we don't care errors here

	// 4. response build UUID
	c.JSON(http.StatusOK, out)
}

func findSSAFunc(code, funcname string) bool {
	// func Foo (
	re := regexp.MustCompile(fmt.Sprintf(`func[ \t]+%s[ \t]*\(`, funcname))
	return re.FindString(code) != ""
}

func isPackageTest(code string) bool {
	// package *_test
	re := regexp.MustCompile(`package .*\_test`)
	return re.FindString(code) != ""
}

func autoimports(outf string) error {
	cmd := exec.Command("goimports", "-w", outf)
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		msg := cmd.Stderr.(*bytes.Buffer).String()
		msg = strings.ReplaceAll(msg, filepath.Dir(outf), "$GOSSAPATH")
		return errors.New(msg)
	}
	return nil
}

func initModules(path string) error {
	// 1. go mod init
	cmd := exec.Command("go", "mod", "init", "gossa")
	cmd.Dir = path
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		msg := cmd.Stderr.(*bytes.Buffer).String()
		msg = strings.ReplaceAll(msg, path, "$GOSSAPATH")
		return errors.New(msg)
	}

	// 2. go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = path
	cmd.Stderr = &bytes.Buffer{}
	err = cmd.Run()
	if err != nil {
		msg := cmd.Stderr.(*bytes.Buffer).String()
		msg = strings.ReplaceAll(msg, path, "$GOSSAPATH")
		return errors.New(msg)
	}

	return nil
}

func buildSSA(funcname, gcflags, outf, buildf string, isTest bool) error {
	var cmd *exec.Cmd
	if !isTest {
		cmd = exec.Command("go", "build", "-mod=readonly", fmt.Sprintf(`-gcflags=%s`, gcflags), "-o", outf, buildf)
	} else {
		cmd = exec.Command("go", "test", "-mod=readonly", fmt.Sprintf(`-gcflags=%s`, gcflags), buildf)
	}
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOSSAFUNC=%s", funcname))
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		msg := cmd.Stderr.(*bytes.Buffer).String()
		msg = strings.ReplaceAll(msg, filepath.Dir(outf), "$GOSSAPATH")
		return errors.New(msg)
	}
	return nil
}
