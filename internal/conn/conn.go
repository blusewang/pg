// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package conn

import (
	"bufio"
	"database/sql/driver"
	"net"
)

type conn struct {
	con       net.Conn
	dsn       dataSourceName
	bufReader *bufio.Reader
	txStatus  TransactionStatus
}

func Open(name string) (c conn, err error) {
	c.dsn, err = ParseDSN(name)
	if err != nil {
		return
	}
	c.con, err = net.DialTimeout(c.dsn.Address())
	if err != nil {
		return
	}
	return
}

func (c conn) Prepare(query string) (stmt driver.Stmt, err error) {
	// TODO
	return
}
func (c conn) Close() (err error) {
	//TODO
	return
}
func (c conn) Begin() (tx driver.Tx, err error) {
	//TODO
	return
}
