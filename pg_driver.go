// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package pg

import (
	"database/sql/driver"
	driver2 "github.com/blusewang/pg/internal/driver"
)

type Driver struct {
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	// 此接口似乎已废弃
	return driver2.NewPgConn(name)
}

func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	return &Connector{Name: name}, nil
}
