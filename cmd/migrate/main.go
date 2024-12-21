package main

import (
	"context"
	"flag"
	"github.com/jiu-u/oai-api/cmd/migrate/wire"
	"github.com/jiu-u/oai-api/pkg/config"
)

func main() {
	var envConf = flag.String("conf", "config/config.yaml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.LoadConfig(*envConf)
	app, cleanup, err := wire.NewWire(conf)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	//logger.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	//logger.Info("docs addr", zap.String("addr", fmt.Sprintf("http://%s:%d/swagger/index.html", conf.GetString("http.host"), conf.GetInt("http.port"))))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
