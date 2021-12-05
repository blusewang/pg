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

type PgStmt struct {
	client.Statement
	client    *client.Client
	resultSig chan int
}

func NewPgStmt(client *client.Client, query string) (st PgStmt, err error) {
	if client.IOError != nil {
		err = driver.ErrBadConn
		return
	}
	var id = fmt.Sprintf("%x", md5.Sum([]byte(query)))

	if client.StatementMaps[id] == nil {
		st.client = client
		st.Id = id
		st.SQL = query
		st.Response, err = client.Parse(st.Id, st.SQL)
		if err == nil {
			client.StatementMaps[id] = &st.Statement
		}
	} else {
		st = PgStmt{client: client, Statement: *client.StatementMaps[id]}
	}
	st.resultSig = make(chan int)
	return
}

func NewPgStmtContext(ctx context.Context, client *client.Client, query string) (st PgStmt, err error) {
	if client.IOError != nil {
		err = driver.ErrBadConn
		return
	}
	var id = fmt.Sprintf("%x", md5.Sum([]byte(query)))

	if client.StatementMaps[id] == nil {
		st.client = client
		st.Id = id
		st.SQL = query
		st.Response, err = client.Parse(st.Id, st.SQL)
		if err == nil {
			client.StatementMaps[id] = &st.Statement
		}
	} else {
		st = PgStmt{client: client, Statement: *client.StatementMaps[id]}
	}
	st.resultSig = make(chan int)
	return
}

func (s PgStmt) Close() (err error) {
	if s.client.IOError != nil {
		return driver.ErrBadConn
	}
	// 当缓存里没有时，关闭语句
	if s.client.StatementMaps[s.Id] == nil {
		return s.client.CloseParse(s.Id)
	}
	return
}

func (s PgStmt) NumInput() int {
	if s.Response.ParameterDescription == nil {
		return 0
	}
	return len(s.Response.ParameterDescription.TypeOIDs)
}

func (s PgStmt) Exec(args []driver.Value) (res driver.Result, err error) {
	if s.client.IOError != nil {
		return nil, driver.ErrBadConn
	}
	response, err := s.client.ParseExec(s.Id, args)
	return driver.RowsAffected(response.Completion.Affected()), err
}

func (s PgStmt) Query(args []driver.Value) (_ driver.Rows, err error) {
	if s.client.IOError != nil {
		return nil, driver.ErrBadConn
	}
	var pr = new(PgRows)
	pr.location = s.client.Location
	pr.columns = s.Response.Description
	res, err := s.client.ParseQuery(s.Id, args)
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
	if s.client.IOError != nil {
		return nil, driver.ErrBadConn
	}
	var as []driver.Value
	for _, v := range args {
		as = append(as, v.Value)
	}
	n, err := s.client.ParseExec(s.Id, as)
	return driver.RowsAffected(n.Completion.Affected()), err
}

// QueryContext executes a query that may return rows, such as a
// SELECT.
//
// QueryContext must honor the context timeout and return when it is canceled.
func (s PgStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (_ driver.Rows, err error) {
	go s.watchCancel(ctx)
	defer s.complete()
	if s.client.IOError != nil {
		return nil, driver.ErrBadConn
	}

	var as []driver.Value
	for _, v := range args {
		as = append(as, v.Value)
	}

	var pr = new(PgRows)
	pr.location = s.client.Location
	pr.columns = s.Response.Description
	res, err := s.client.ParseQuery(s.Id, as)
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
	_ = s.client.CancelRequest()
}

func (s PgStmt) complete() {
	s.resultSig <- 1
}
