package app

import (
	"context"
	"database/sql/driver"
	"github.com/blusewang/pg/v2/internal/client"
)

type Connector struct {
	dsn    client.DataSourceName
	driver Driver
}

func (c Connector) Connect(ctx context.Context) (conn driver.Conn, err error) {
	return NewConnect(ctx, c.dsn)
}

func (c Connector) Driver() driver.Driver {
	return c.driver
}

var _ driver.Connector = new(Connector)
