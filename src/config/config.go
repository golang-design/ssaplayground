package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		logrus.Fatalf("fatal: fail to read configuration file: %v", err)
	}

	conf = &Config{}
	err = yaml.Unmarshal(y, conf)
	if err != nil {
		logrus.Fatalf("fatal: fail to parse configuration file: %v", err)
	}
	gin.SetMode(conf.Mode)

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetReportCaller(false)
	logrus.Infof("load config file: %q", f)
}
