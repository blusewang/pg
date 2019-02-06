// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"database/sql/driver"
	"reflect"
)

type PgRows struct {
}

func (pr *PgRows) Columns() []string {
	return nil
}

func (pr *PgRows) Close() error {
	return nil
}

func (pr *PgRows) Next(dest []driver.Value) error {
	return nil
}

// may be implemented by Rows. It should return the precision and scale for decimal types.
// If not applicable, ok should be false.
func (pr *PgRows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	return 0, 0, false
}

// may be implemented by Rows. The nullable value should be true if it is known the column may be null,
// or false if the column is known to be not nullable. If the column nullability is unknown, ok should be false.
func (pr *PgRows) ColumnTypeNullable(index int) (nullable, ok bool) {
	return false, false
}

func (pr *PgRows) ColumnTypeLength(index int) (length int64, ok bool) {
	return 0, false
}

func (pr *PgRows) ColumnTypeDatabaseTypeName(index int) string {
	return ""
}

func (pr *PgRows) ColumnTypeScanType(index int) reflect.Type {
	return nil
}

// RowsNextResultSet
// HasNextResultSet is called at the end of the current result set and
// reports whether there is another result set after the current one.
func (pr *PgRows) HasNextResultSet() bool {
	return false
}

// NextResultSet advances the driver to the next result set even
// if there are remaining rows in the current result set.
//
// NextResultSet should return io.EOF when there are no more result sets.
func (pr *PgRows) NextResultSet() error {
	return nil
}
