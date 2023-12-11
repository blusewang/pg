// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package pg

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/blusewang/pg/client"
	"github.com/blusewang/pg/dsn"
	"log"
)

const DriverName = "pg"

func init() {
	log.Println("register driver")
	sql.Register(DriverName, &Driver{})
}

func NewConnector(dataSourceName string) driver.Connector {
	return &Connector{Name: dataSourceName}
}

func New(ctx context.Context, dsn dsn.DataSourceName) (*client.Client, error) {
	return client.NewClient(ctx, dsn)
}
