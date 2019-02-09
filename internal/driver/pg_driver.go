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

type PgDriver struct{}

func (d *PgDriver) Open(name string) (driver.Conn, error) {
	return NewPgConn(name)
}

// BeginTx starts and returns a new transaction.
// If the context is canceled by the user the sql package will
// call Tx.Rollback before discarding and closing the connection.
//
// This must check opts.Isolation to determine if there is a set
// isolation level. If the driver does not support a non-default
// level and one is set or if there is a non-default isolation level
// that is not supported, an error must be returned.
//
// This must also check opts.ReadOnly to determine if the read-only
// value is true to either set the read-only transaction property if supported
// or return an error if it is not supported.
func (d *PgDriver) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return nil, nil
}
