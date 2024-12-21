package server

import (
	"context"
	"fmt"
	"github.com/jiu-u/oai-api/internal/model"
	"gorm.io/gorm"
	"log"
	"os"
)

type Migrate struct {
	db  *gorm.DB
	log *log.Logger
}

func NewMigrate(db *gorm.DB) *Migrate {
	return &Migrate{
		db: db,
	}
}
func (m *Migrate) Start(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		new(model.Model),
		new(model.Provider),
	); err != nil {
		fmt.Println("AutoMigrate error", err)
		return err
	}
	fmt.Println("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *Migrate) Stop(ctx context.Context) error {
	fmt.Println("AutoMigrate stop")
	return nil
}
