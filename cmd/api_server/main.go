package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jiu-u/oai-api/cmd/api_server/wire"
	"github.com/jiu-u/oai-api/cmd/api_server/wire_load"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/log"
	"go.uber.org/zap"
)

func main() {
	var envConf = flag.String("conf", "config/config.yaml", "config path, eg: -conf ./config/local.yml")
	var load = flag.Bool("load", false, "load data from yaml file,eg: -load true")
	flag.Parse()
	conf := config.LoadConfig(*envConf)
	logger := log.NewLogger(conf)
	if *load {
		LoadDataFromFile(conf, logger)
	}
	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info(
		"server start",
		zap.String("host", fmt.Sprintf("http://%s:%d", conf.HTTP.Host, conf.HTTP.Port)),
		zap.String("type", "system"),
	)
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}

func LoadDataFromFile(cfg *config.Config, logger *log.Logger) {
	server, cleanup, err := wire_load.NewWire(cfg, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	if err = server.Start(context.Background()); err != nil {
		fmt.Println("load data error", err)
		panic(err)
	}
}
