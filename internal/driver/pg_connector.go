// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"context"
	"database/sql/driver"
)

type PgConnector struct {
	Name string
}

func (c *PgConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return NewPgConnContext(ctx, c.Name)
}

func (c *PgConnector) Driver() driver.Driver {
	return &PgDriver{}
}
