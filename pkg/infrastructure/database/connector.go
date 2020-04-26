package database

import "database/sql"

type Connector interface {
	Client() Client
	Close() error
}

type Client interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
