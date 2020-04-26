package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const driverName = "mysql"

type mysqlConnector struct {
	db *sql.DB
}

func (c *mysqlConnector) Client() Client {
	return c.db
}

func (c *mysqlConnector) Close() error {
	return c.db.Close()
}

func NewMySQLConnector(dsn *DSN) (Connector, error) {
	db, err := sql.Open(driverName, dsn.ToString())
	if err != nil {
		return nil, err
	}
	return &mysqlConnector{db: db}, nil
}
