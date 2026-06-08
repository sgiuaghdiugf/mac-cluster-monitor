package database

import (
	"database/sql"
	"mac-cluster-monitor/server/config"
	
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init(cfg *config.Config) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("mysql", cfg.Database.DSN())
	if err != nil {
		return nil, err
	}
	
	if err = DB.Ping(); err != nil {
		return nil, err
	}
	
	DB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	
	return DB, nil
}
