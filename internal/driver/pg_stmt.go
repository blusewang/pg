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

type PgStmt struct{}

func (s *PgStmt) Close() error {
	return nil
}
func (s *PgStmt) NumInput() int {
	return 0
}

func (s *PgStmt) Exec(args []driver.Value) (driver.Result, error) {
	return &PgResult{}, nil
}

func (s *PgStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &PgRows{}, nil
}

// ExecContext executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// ExecContext must honor the context timeout and return when it is canceled.
func (s *PgStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return &PgResult{}, nil
}

// QueryContext executes a query that may return rows, such as a
// SELECT.
//
// QueryContext must honor the context timeout and return when it is canceled.
func (s *PgStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return &PgRows{}, nil
}
