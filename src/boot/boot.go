package boot

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/changkun/gossafunc/src/config"
	"github.com/changkun/gossafunc/src/route"
	"github.com/sirupsen/logrus"
)

func init() {
	config.Init()
}

func Run() {
	server := &http.Server{
		Handler: route.Register(),
		Addr:    config.Get().Addr,
	}

	terminated := make(chan bool, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill)
		sig := <-quit

		logrus.Info("gossaweb: service is stopped with signal: ", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			logrus.Errorf("gossaweb: close gossaweb with error: %v", err)
		}

		cancel()
		terminated <- true
	}()

	logrus.Infof("gossaweb: welcome to gossaweb service... http://%s/gossa", config.Get().Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		logrus.Info("gossaweb: launch with error: ", err)
	}

	<-terminated
	logrus.Info("gossaweb: service has terminated successfully, good bye!")
}
