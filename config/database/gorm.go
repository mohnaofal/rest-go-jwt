package database

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type gormConnector struct {
	gormMysql *gorm.DB
}

type GormConnector interface {
	MysqlGorm() *gorm.DB
}

func (c *gormConnector) MysqlGorm() *gorm.DB {
	return c.gormMysql
}

func ConnectorDB() GormConnector {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := os.Getenv("MYSQL_HOST")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &gormConnector{
		gormMysql: db,
	}
}
