// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package pg

import (
	"database/sql/driver"
	"github.com/blusewang/pg/internal/conn"
)

type Driver struct{}

func (d *Driver) Open(name string) (_ driver.Conn, err error) {

	c, err := conn.Open(name)

	return c, err
}
