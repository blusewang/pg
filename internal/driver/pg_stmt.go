// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/blusewang/pg/internal/network"
	"hash/crc32"
	"log"
)

func NewPgStmt(io *network.PgIO, query string, dsn *DataSourceName) (st *PgStmt, err error) {
	if io.IOError != nil {
		return nil, driver.ErrBadConn
	}
	st = new(PgStmt)
	st.io = io
	st.dsn = dsn
	st.Identifies = fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(query)))
	st.Sql = query
	st.columns, st.parameterTypes, err = st.io.Parse(st.Identifies, st.Sql)
	st.resultSig = make(chan int)
	return st, err
}

type PgStmt struct {
	io             *network.PgIO
	dsn            *DataSourceName
	Identifies     string
	Sql            string
	columns        []network.PgColumn
	parameterTypes []uint32
	resultSig      chan int
}

func (s *PgStmt) Close() error {
	if s.io.IOError != nil {
		return driver.ErrBadConn
	}
	close(s.resultSig)
	return s.io.CloseParse(s.Identifies)
}

func (s *PgStmt) NumInput() int {
	return len(s.parameterTypes)
}

func (s *PgStmt) Exec(args []driver.Value) (res driver.Result, err error) {
	log.Println("exec")
	if s.io.IOError != nil {
		return nil, driver.ErrBadConn
	}
	var as []interface{}
	for _, v := range args {
		as = append(as, v)
	}
	n, err := s.io.ParseExec(s.Identifies, as)
	return driver.RowsAffected(n), err
}

func (s *PgStmt) Query(args []driver.Value) (_ driver.Rows, err error) {
	var as []interface{}
	for _, v := range args {
		as = append(as, v)
	}

	var pr = new(PgRows)
	pr.position = -1
	pr.columns = s.columns
	pr.parameterTypes = s.parameterTypes
	pr.fieldLen, pr.rows, err = s.io.ParseQuery(s.Identifies, as)
	return pr, nil
}

// ExecContext executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// ExecContext must honor the context timeout and return when it is canceled.
func (s *PgStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	go s.watchCancel(ctx)
	defer s.complete()
	if s.io.IOError != nil {
		return nil, driver.ErrBadConn
	}
	var as []interface{}
	for _, v := range args {
		as = append(as, v.Value)
	}
	n, err := s.io.ParseExec(s.Identifies, as)
	return driver.RowsAffected(n), err
}

// QueryContext executes a query that may return rows, such as a
// SELECT.
//
// QueryContext must honor the context timeout and return when it is canceled.
func (s *PgStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (_ driver.Rows, err error) {
	go s.watchCancel(ctx)
	defer s.complete()
	if s.io.IOError != nil {
		return nil, driver.ErrBadConn
	}

	var as []interface{}
	for _, v := range args {
		as = append(as, v.Value)
	}

	var pr = new(PgRows)
	pr.location = s.io.Location
	pr.columns = s.columns
	pr.parameterTypes = s.parameterTypes
	pr.fieldLen, pr.rows, err = s.io.ParseQuery(s.Identifies, as)

	return pr, nil
}

func (s *PgStmt) watchCancel(ctx context.Context) {
	select {
	case <-ctx.Done():
		s.cancel()
	case <-s.resultSig:
	}
}

func (s *PgStmt) cancel() {
	_ = s.io.CancelRequest(s.dsn.Address())
}

func (s *PgStmt) complete() {
	s.resultSig <- 1
}
