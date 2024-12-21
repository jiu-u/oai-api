//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/server"
	"github.com/jiu-u/oai-api/pkg/app"
	"github.com/jiu-u/oai-api/pkg/config"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewTransaction,
)

var serverSet = wire.NewSet(
	server.NewMigrate,
)

// build App
func newApp(
	migrateServer *server.Migrate,
	// job *server.Job,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(migrateServer),
		app.WithName("demo-server"),
	)
}

func NewWire(cfg *config.Config) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serverSet,
		newApp,
	))
}
