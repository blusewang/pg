// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/blusewang/pg/internal/frame"
	"io"
	"math"
	"reflect"
	"time"
)

const headerSize = 4

type PgRows struct {
	isStrict bool
	location *time.Location
	columns  *frame.RowDescription
	rows     []*frame.DataRow
	position int
}

func (pr *PgRows) Columns() (cols []string) {
	for _, v := range pr.columns.Columns {
		cols = append(cols, v.Name)
	}
	return
}

func (pr *PgRows) Close() error {
	pr.position = 0
	pr.rows = nil
	pr.columns = nil
	return nil
}

func (pr *PgRows) Next(dest []driver.Value) error {
	var rowsLen = len(pr.rows)
	if rowsLen == 0 {
		return sql.ErrNoRows
	} else if pr.position < 0 || pr.position >= rowsLen {
		return fmt.Errorf("pg_rows rows length is %v but position is %v", rowsLen, pr.position)
	} else if pr.position == rowsLen {
		return io.EOF
	}
	for k, v := range (pr.rows)[pr.position].DataArr {
		dest[k] = convert(v, pr.columns.Columns[k], pr.location, pr.isStrict)
	}
	pr.position += 1
	return nil
}

// ColumnTypePrecisionScale may be implemented by Rows. It should return the precision and scale for decimal types.
// If not applicable, ok should be false.
func (pr *PgRows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	var fd = pr.columns.Columns[index]
	switch fd.TypeOid {
	case PgTypeNumeric, PgTypeArrNumeric:
		mod := fd.TypeModifier - headerSize
		precision = int64((mod >> 16) & 0xffff)
		scale = int64(mod & 0xffff)
		return precision, scale, true
	default:
		return 0, 0, false
	}
}

func (pr *PgRows) ColumnTypeLength(index int) (length int64, ok bool) {
	switch pr.columns.Columns[index].TypeOid {
	case PgTypeText, PgTypeBytea:
		return math.MaxInt64, true
	case PgTypeVarchar, PgTypeBpchar:
		return int64(pr.columns.Columns[index].TypeModifier - headerSize), true
	default:
		return 0, false
	}
}

func (pr *PgRows) ColumnTypeDatabaseTypeName(index int) string {
	return PgTypeMap[PgType(pr.columns.Columns[index].TypeOid)]
}

func (pr *PgRows) ColumnTypeScanType(index int) reflect.Type {
	switch PgType(pr.columns.Columns[index].TypeOid) {
	case PgTypeBool:
		return reflect.TypeOf(false)
	case PgTypeDate, PgTypeTime, PgTypeTimestamp, PgTypeTimestamptz, PgTypeTimetz:
		return reflect.TypeOf(time.Time{})
	case PgTypeInt2, PgTypeInt4, PgTypeInt8:
		return reflect.TypeOf(int64(0))
	case PgTypeFloat4, PgTypeFloat8, PgTypeNumeric:
		return reflect.TypeOf(float64(0))
	case PgTypeText, PgTypeVarchar, PgTypeChar, PgTypeUuid, PgTypeJson, PgTypeJsonb, PgTypePoint:
		return reflect.TypeOf("")
	case PgTypeBytea:
		return reflect.TypeOf([]byte{})

	case PgTypeArrBool:
		return reflect.TypeOf([]bool{})
	case PgTypeArrDate, PgTypeArrTime, PgTypeArrTimestamp, PgTypeArrTimestamptz, PgTypeArrTimetz:
		return reflect.TypeOf([]time.Time{})
	case PgTypeArrInt2, PgTypeArrInt4, PgTypeArrInt8:
		return reflect.TypeOf([]int64{})
	case PgTypeArrFloat4, PgTypeArrFloat8, PgTypeArrNumeric:
		return reflect.TypeOf([]float64{})
	case PgTypeArrVarchar, PgTypeArrChar, PgTypeArrText, PgTypeArrUuid, PgTypeArrJson, PgTypeArrJsonb:
		return reflect.TypeOf([]string{})
	case PgTypeArrBytea:
		return reflect.TypeOf([][]byte{})

	default:
		return reflect.TypeOf("")
	}
}

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
