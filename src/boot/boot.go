// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a GPLv3 license that can be found in the LICENSE file.

package boot

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.design/x/ssaplayground/src/config"
	"golang.design/x/ssaplayground/src/route"
)

func init() {
	log.SetPrefix("redir: ")
	log.SetFlags(log.Lmsgprefix | log.LstdFlags | log.Lshortfile)
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
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		sig := <-quit

		log.Printf("service is stopped with signal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("close ssaplayground with error: %v", err)
		}

		cancel()
		terminated <- true
	}()

	log.Printf("welcome to ssaplayground service... http://%s/gossa", config.Get().Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		terminated <- true
		log.Printf("launch with error: %v", err)
	}

	<-terminated
	log.Printf("service has terminated successfully, good bye!")
}
