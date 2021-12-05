// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"database/sql/driver"
	"errors"
)

type PgTx struct {
	pgConn PgConn
}

func (t PgTx) Commit() (err error) {
	if t.pgConn.client.IOError != nil {
		return driver.ErrBadConn
	}
	if t.pgConn.client.IsInTransaction() == false {
		err = errors.New("this connection is out of transaction")
	}
	_, err = t.pgConn.client.QueryNoArgs("commit")
	if err != nil {
		return
	}
	if t.pgConn.client.IsInTransaction() {
		err = errors.New("commit fail")
	}
	return
}

func (t PgTx) Rollback() (err error) {
	if t.pgConn.client.IOError != nil {
		return driver.ErrBadConn
	}
	if t.pgConn.client.IsInTransaction() == false {
		err = errors.New("this connection is out of transaction")
	}
	_, err = t.pgConn.client.QueryNoArgs("rollback")
	if err != nil {
		return
	}
	if t.pgConn.client.IsInTransaction() {
		err = errors.New("rollback fail")
	}
	return nil
}
