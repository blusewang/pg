// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"context"
	"crypto/md5"
	"database/sql/driver"
	"fmt"
	"github.com/blusewang/pg/internal/client"
)

func NewPgStmt(conn PgConn, query string) (st PgStmt, err error) {
	if conn.io.IOError != nil {
		err = conn.io.Err.Error
		return
	}
	var id = fmt.Sprintf("%x", md5.Sum([]byte(query)))
	st = conn.stmts[id]
	if st.Identifies == "" {
		st = *new(PgStmt)
		st.pgConn = conn
		st.Identifies = id
		st.Sql = query
		st.Response, err = st.pgConn.io.Parse(st.Identifies, st.Sql)
		st.resultSig = make(chan int)
		conn.stmts[id] = st
	}
	return st, err
}

type PgStmt struct {
	pgConn     PgConn
	Identifies string
	Sql        string
	Response   client.Response
	resultSig  chan int
}

func (s PgStmt) Close() (err error) {
	if s.pgConn.io.IOError != nil {
		return s.pgConn.io.Err.Error
	}
	err = s.pgConn.io.CloseParse(s.Identifies)
	if err != nil {
		return
	}
	close(s.resultSig)
	if s.pgConn.stmts[s.Identifies].Identifies != "" {
		delete(s.pgConn.stmts, s.Identifies)
	}
	return nil
}

func (s PgStmt) NumInput() int {
	if s.Response.Description == nil {
		return 0
	}
	return len(s.Response.Description.Columns)
}

func (s PgStmt) Exec(args []driver.Value) (res driver.Result, err error) {
	if s.pgConn.io.IOError != nil {
		return nil, s.pgConn.io.Err.Error
	}
	var as [][]byte
	for _, v := range args {
		as = append(as, driverValue2Pg(v))
	}
	response, err := s.pgConn.io.ParseExec(s.Identifies, as)
	return driver.RowsAffected(response.Completion.Affected()), err
}

func (s PgStmt) Query(args []driver.Value) (_ driver.Rows, err error) {
	var as [][]byte
	for _, v := range args {
		as = append(as, driverValue2Pg(v))
	}

	var pr = new(PgRows)
	pr.isStrict = s.pgConn.dsn.IsStrict
	pr.location = s.pgConn.io.Location
	pr.columns = s.Response.Description
	res, err := s.pgConn.io.ParseQuery(s.Identifies, as)
	pr.rows = res.DataRows
	return pr, err
}

// ExecContext executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// ExecContext must honor the context timeout and return when it is canceled.
func (s PgStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	go s.watchCancel(ctx)
	defer s.complete()
	if s.pgConn.io.IOError != nil {
		return nil, s.pgConn.io.Err.Error
	}
	var as [][]byte
	for _, v := range args {
		as = append(as, driverValue2Pg(v.Value))
	}
	n, err := s.pgConn.io.ParseExec(s.Identifies, as)
	return driver.RowsAffected(n.Completion.Affected()), err
}

// QueryContext executes a query that may return rows, such as a
// SELECT.
//
// QueryContext must honor the context timeout and return when it is canceled.
func (s PgStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (_ driver.Rows, err error) {
	go s.watchCancel(ctx)
	defer s.complete()
	if s.pgConn.io.IOError != nil {
		return nil, s.pgConn.io.Err.Error
	}

	var as [][]byte
	for _, v := range args {
		as = append(as, driverValue2Pg(v.Value))
	}

	var pr = new(PgRows)
	pr.isStrict = s.pgConn.dsn.IsStrict
	pr.location = s.pgConn.io.Location
	pr.columns = s.Response.Description
	res, err := s.pgConn.io.ParseQuery(s.Identifies, as)
	pr.rows = res.DataRows
	return pr, err
}

func (s PgStmt) watchCancel(ctx context.Context) {
	select {
	case <-ctx.Done():
		s.cancel()
	case <-s.resultSig:
	}
}

func (s PgStmt) cancel() {
	_ = s.pgConn.io.CancelRequest()
}

func (s PgStmt) complete() {
	s.resultSig <- 1
}
