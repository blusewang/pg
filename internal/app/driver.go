package app

import (
	"context"
	"database/sql/driver"
	"github.com/blusewang/pg/v2/internal/client"
)

type Driver struct {
}

func (d Driver) Open(name string) (c driver.Conn, err error) {
	dsn, err := client.ParseDSN(name)
	if err != nil {
		return
	}
	return NewConnect(context.Background(), dsn)
}

func (d Driver) OpenConnector(name string) (driver.Connector, error) {
	dsn, err := client.ParseDSN(name)
	return Connector{dsn: dsn, driver: d}, err
}

var _ driver.Driver = new(Driver)
var _ driver.DriverContext = new(Driver)
