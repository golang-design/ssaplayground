// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a GPLv3 license that can be found in the LICENSE file.

package boot

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"golang.design/x/ssaplayground/src/config"
	"golang.design/x/ssaplayground/src/route"
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

		logrus.Info("ssaplayground: service is stopped with signal: ", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			logrus.Errorf("ssaplayground: close ssaplayground with error: %v", err)
		}

		cancel()
		terminated <- true
	}()

	logrus.Infof("ssaplayground: welcome to ssaplayground service... http://%s/gossa", config.Get().Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		logrus.Info("ssaplayground: launch with error: ", err)
	}

	<-terminated
	logrus.Info("ssaplayground: service has terminated successfully, good bye!")
}
