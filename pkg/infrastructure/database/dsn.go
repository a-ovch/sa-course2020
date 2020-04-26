package database

import "fmt"

type DSN struct {
	host string
	port string
	user string
	pwd  string
	db   string
}

func (d *DSN) ToString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.user, d.pwd, d.host, d.port, d.db)
}

func NewDSN(host, port, user, pwd, db string) *DSN {
	return &DSN{host: host, port: port, user: user, pwd: pwd, db: db}
}
