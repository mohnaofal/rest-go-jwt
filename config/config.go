package config

import "github.com/mohnaofal/rest-go-jwt/config/database"

type Config struct {
	db database.GormConnector
}

func InitConfig() *Config {
	cfg := new(Config)

	cfg.InitDB()

	return cfg
}

func (c *Config) InitDB() {
	c.db = database.ConnectorDB()
}

func (c *Config) DB() database.GormConnector {
	return c.db
}
