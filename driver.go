// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package pg

import (
	"database/sql/driver"
	"fmt"
	"github.com/blusewang/pg/internal/conn"
)

type Driver struct{}

func (d *Driver) Open(name string) (_ driver.Conn, err error) {
	//
	defer func(err *error) {
		if e := recover(); e != nil {
			*err = fmt.Errorf("pg: Open Error: %#v", e)
		}
	}(&err)

	c, err := conn.Open(name)

	return c, err
}
