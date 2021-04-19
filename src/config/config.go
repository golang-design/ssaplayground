// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a GPLv3 license that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr   string `json:"addr"`
	Mode   string `json:"mode"`
	Static string `json:"static"`
}

var conf *Config

func Get() *Config {
	return conf
}

func Init() {
	c := flag.String("conf", "", "path to the ssaplayground config file")
	usage := func() {
		fmt.Fprintf(os.Stderr, `
SSAPLAYGROUND is a web service for exploring Go's SSA intermediate representation.
Usage:
`)
		flag.PrintDefaults()
	}
	flag.Usage = usage
	flag.Parse()
	f := *c
	if len(f) == 0 {
		f = os.Getenv("GOSSAWEB_CONF")
	}
	if len(f) == 0 {
		usage()
		os.Exit(1)
	}

	y, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("fatal: fail to read configuration file: %v", err)
	}

	conf = &Config{}
	err = yaml.Unmarshal(y, conf)
	if err != nil {
		log.Fatalf("fatal: fail to parse configuration file: %v", err)
	}
	gin.SetMode(conf.Mode)

	log.Printf("load config file: %q", f)
}
