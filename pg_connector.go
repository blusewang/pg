// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package pg

import (
	"context"
	"database/sql/driver"
	driver2 "github.com/blusewang/pg/internal/driver"
)

type Connector struct {
	Name string
}

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return driver2.NewPgConnContext(ctx, c.Name)
}

func (c *Connector) Driver() driver.Driver {
	return &Driver{}
}
