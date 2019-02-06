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

type PgConnector struct{}

func (c *PgConnector) Connect(context.Context) (driver.Conn, error) {
	return &pgConn{}, nil
}

func (c *PgConnector) Driver() driver.Driver {
	return &PgDriver{}
}
