// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a GPLv3 license that can be found in the LICENSE file.

package route

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.design/x/ssaplayground/src/config"
)

// Register routers for the ssa service.
func Register() *gin.Engine {
	r := &Router{Engine: gin.Default()}
	r.SetupAPI()
	r.SetupApp()
	if config.Get().Mode != gin.DebugMode {
		return r.Engine
	}
	r.SetupProfile()
	return r.Engine
}

// Router is a router engine.
type Router struct {
	Engine *gin.Engine
}

// SetupAPI serves the API endpoints of the gossa service.
func (r *Router) SetupAPI() {
	v1 := r.Engine.Group("/gossa/api/v1")
	{
		v1.GET("/ping", Pong)
		v1.POST("/buildssa", BuildSSA)
	}
}

// SetupApp serves the static website of Go SSA Playground.
func (r *Router) SetupApp() {
	r.Engine.Use(static("/gossa"))
	log.Printf("GoSSAWeb is on: http://%s, static: %s", config.Get().Addr, config.Get().Static)
}

// SetupProfile the standard HandlerFuncs from the net/http/pprof package with
// the provided gin.Engine. prefixOptions is a optional. If not prefixOptions,
// the default path prefix is used, otherwise first prefixOptions will be path prefix.
//
// Basic Usage:
//
// - use the pprof tool to look at the heap profile:
//   go tool pprof localhost:8080/midgard/api/v1/debug/pprof/heap
// - look at a 30-second CPU profile:
//   go tool pprof localhost:8080/midgard/api/v1/debug/pprof/profile
// - look at the goroutine blocking profile, after calling runtime.SetBlockProfileRate:
//   go tool pprof localhost:8080/midgard/api/v1/debug/pprof/block
// - collect a 5-second execution trace:
//   go tool pprof localhost:8080/midgard/api/v1/debug/pprof/trace?seconds=5
//
func (r *Router) SetupProfile() {
	pprofHandler := func(h http.HandlerFunc) gin.HandlerFunc {
		handler := http.HandlerFunc(h)
		return func(c *gin.Context) {

			fmt.Println(c.Request.Host)
			if !strings.Contains(c.Request.Host, "localhost") {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			handler.ServeHTTP(c.Writer, c.Request)
		}
	}
	rr := r.Engine.Group("/debug/pprof")
	{
		rr.GET("/", pprofHandler(pprof.Index))
		rr.GET("/cmdline", pprofHandler(pprof.Cmdline))
		rr.GET("/profile", pprofHandler(pprof.Profile))
		rr.POST("/symbol", pprofHandler(pprof.Symbol))
		rr.GET("/symbol", pprofHandler(pprof.Symbol))
		rr.GET("/trace", pprofHandler(pprof.Trace))
		rr.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		rr.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		rr.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		rr.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		rr.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		rr.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	}
}
