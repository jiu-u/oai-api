package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jiu-u/oai-api/cmd/server/wire"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/log"
	"go.uber.org/zap"
)

func main() {
	//time.Sleep(5 * time.Second)
	var envConf = flag.String("conf", "config/local.yaml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.LoadConfig(*envConf)
	logger := log.NewLogger(conf)
	app, cleanup, err := wire.NewWire(conf, logger)
	if err != nil {
		logger.Error("wire error", zap.Error(err))
		panic(err)
	}
	defer cleanup()
	// 创建表
	migrate := app.MigrateJob
	err = migrate.Start(context.Background())
	if err != nil {
		logger.Error("migrate error", zap.Error(err))
		panic(err)
	}
	// 读取yaml文件
	apiApp := app.App
	logger.Info(
		"server start",
		zap.String("host", fmt.Sprintf("http://%s:%d", conf.HTTP.Host, conf.HTTP.Port)),
		zap.String("type", "system"),
	)
	if err = apiApp.Run(context.Background()); err != nil {
		panic(err)
	}
}
