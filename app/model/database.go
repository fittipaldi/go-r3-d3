package model

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/fittipaldi/go-r3-d3/config"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DatabaseConnection(dbCredential *config.DBConfig) *sql.DB {
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbCredential.Username, dbCredential.Password, dbCredential.Host, dbCredential.Port, dbCredential.Name)
	db, err := sql.Open(dbCredential.Type, dbUri)

	if err != nil {
		fmt.Sprintf("ERROR TO CONNECT TO DATABASE [%s]", dbCredential.Host)
		log.Fatal(err)
	}
	return db
}

func InitGorm(dbType string, config *config.Config, db *sql.DB) *gorm.DB {
	return initGormMySQL(config, db)
}

func initGormMySQL(config *config.Config, db *sql.DB) *gorm.DB {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Info),
		CreateBatchSize: config.DB.DBBatchSize,
	})
	if err != nil {
		fmt.Sprintf("ERROR ON GORM [%s]", err.Error())
		log.Fatal("Failed to initialize *gorm.DB with an existing database connection")
	}
	return gormDB
}