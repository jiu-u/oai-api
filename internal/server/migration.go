package server

import (
	"context"
	"fmt"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Migrate struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewMigrate(db *gorm.DB, logger *log.Logger) *Migrate {
	return &Migrate{
		db:     db,
		logger: logger,
	}
}
func (m *Migrate) Start(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		new(model.ChannelModel),
		new(model.Channel),
		//new(model.Model),
		//new(model.Provider),
		new(model.User),
		new(model.ApiKey),
		new(model.RequestLog),
		new(model.SystemConfig),
		new(model.AsyncTask),
		new(model.UserAuthProvider),
	); err != nil {
		m.logger.Error("AutoMigrate error", zap.Error(err))
		return err
	}
	m.logger.Info("AutoMigrate success")
	//os.Exit(0)
	return nil
}
func (m *Migrate) Stop(ctx context.Context) error {
	fmt.Println("AutoMigrate stop")
	return nil
}
