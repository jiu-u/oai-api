package repository

import (
	"github.com/glebarez/sqlite"
	"github.com/jiu-u/oai-api/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"os"

	//"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

func NewDB(cfg *config.Config) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	//logger := zapgorm2.New(l.Logger)
	driver := cfg.Database.Driver
	dsn := cfg.Database.Dsn

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html
	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	case "sqlite":
		_, err := os.Stat("./data/db/oai.db")
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll("./data/db", os.ModePerm)
				if err != nil {
					panic(err)
				}
				f, err := os.Create("./data/db/oai.db")
				if err != nil {
					panic(err)
				}
				f.Close()
			}
		}
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		panic("unknown db driver")
	}
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
