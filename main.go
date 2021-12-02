// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package pg

import (
	"database/sql"
	"database/sql/driver"
	"github.com/blusewang/pg/internal/client"
	dr "github.com/blusewang/pg/internal/driver"
)

func init() {
	sql.Register("pg", &dr.PgDriver{})
}

func NewConnector(dataSourceName string) driver.Connector {
	return &dr.PgConnector{Name: dataSourceName}
}

func Listen(channel string, handler func(string)) {
	client.ListenMap[channel] = handler
}
