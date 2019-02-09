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

func NewPgStmt(io *network.PgIO, query string) (st *PgStmt, err error) {
	if io.IOError != nil {
		return nil, driver.ErrBadConn
	}
	st = new(PgStmt)
	st.io = io
	st.Identifies = fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(query)))
	st.Sql = query
	st.columns, st.parameterTypes, err = st.io.Parse(st.Identifies, st.Sql)
	return st, err
}

type PgStmt struct {
	io             *network.PgIO
	Identifies     string
	Sql            string
	columns        []network.PgColumn
	parameterTypes []uint32
}

func (s *PgStmt) Close() error {
	if s.io.IOError != nil {
		return driver.ErrBadConn
	}
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
	data, err := s.io.ParseExec(s.Identifies, args)
	if err != nil {
		return
	}
	log.Println(data)
	return &PgResult{}, nil
}

func (s *PgStmt) Query(args []driver.Value) (driver.Rows, error) {
	log.Println("query")
	return &PgRows{}, nil
}

// ExecContext executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// ExecContext must honor the context timeout and return when it is canceled.
func (s *PgStmt) ExecContextt(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	log.Println("ExecContext")
	return &PgResult{}, nil
}

// QueryContext executes a query that may return rows, such as a
// SELECT.
//
// QueryContext must honor the context timeout and return when it is canceled.
func (s *PgStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return &PgRows{}, nil
}
