package route

import (
	"net/http"
	"net/http/pprof"

	"github.com/changkun/gossafunc/src/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Register routers
func Register() *gin.Engine {
	r := &Router{
		Engine: gin.Default(),
	}
	r.SetupAPI()
	r.SetupApp()
	if config.Get().Mode != "debug" {
		return r.Engine
	}
	r.SetupProfile()
	return r.Engine
}

type Router struct {
	Engine *gin.Engine
}

func (r *Router) SetupAPI() {
	v1 := r.Engine.Group("/api/v1")
	{
		v1.GET("/ping", Pong)
		v1.POST("/buildssa", BuildSSA)
	}
}

func (r *Router) SetupApp() {
	r.Engine.Use(static("/gossa"))
	logrus.Infof("GoSSAWeb is on: http://%s, static: %s", config.Get().Addr, config.Get().Static)
}

// profile the standard HandlerFuncs from the net/http/pprof package with
// the provided gin.Engine. prefixOptions is a optional. If not prefixOptions,
// the default path prefix is used, otherwise first prefixOptions will be path prefix.
//
// Basic Usage:
//
// - use the pprof tool to look at the heap profile:
//   go tool pprof http://localhost:9999/debug/pprof/heap
// - look at a 30-second CPU profile:
//   go tool pprof http://localhost:9999/debug/pprof/profile
// - look at the goroutine blocking profile, after calling runtime.SetBlockProfileRate:
//   go tool pprof http://localhost:9999/debug/pprof/block
// - collect a 5-second execution trace:
//   wget http://localhost:9999/debug/pprof/trace?seconds=5
//
func (r *Router) SetupProfile() {
	pprofHandler := func(h http.HandlerFunc) gin.HandlerFunc {
		handler := http.HandlerFunc(h)
		return func(c *gin.Context) {
			handler.ServeHTTP(c.Writer, c.Request)
		}
	}
	prefixRouter := r.Engine.Group("/debug/pprof")
	{
		prefixRouter.GET("/", pprofHandler(pprof.Index))
		prefixRouter.GET("/cmdline", pprofHandler(pprof.Cmdline))
		prefixRouter.GET("/profile", pprofHandler(pprof.Profile))
		prefixRouter.POST("/symbol", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/symbol", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/trace", pprofHandler(pprof.Trace))
		prefixRouter.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		prefixRouter.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		prefixRouter.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	}
}
