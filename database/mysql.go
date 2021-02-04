package database

import (
    "database/sql"
    "fmt"
    "github.com/cemalkilic/shorten-backend/config"
    "log"
)

func NewMySQLDBHandler(cfg *config.Config) *sql.DB {

    mysqlUsername := cfg.MysqlUser
    mysqlPassword := cfg.MysqlPass
    mysqlDBName := cfg.MysqlDb
    mysqlPort := cfg.MysqlPort
    mysqlHost := cfg.MysqlHost

    connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)

    database, err := sql.Open("mysql", connStr)
    if err != nil {
        log.Fatal(err)
    }

    return database
}
