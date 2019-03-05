// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"database/sql/driver"
	"fmt"
	"github.com/blusewang/pg/internal/network"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func NewPgConn(name string) (c *pgConn, err error) {
	c = new(pgConn)
	c.dsn, err = ParseDSN(name)
	if err != nil {
		return
	}

	c.io = network.NewPgIO()
	err = c.io.Dial(c.dsn.Address())
	if err != nil {
		return
	}
	err = c.io.StartUp(c.dsn.Parameter, c.dsn.password)
	if err != nil {
		return
	}
	return
}

type pgConn struct {
	dsn *DataSourceName
	io  *network.PgIO
}

// Prepare returns a prepared statement, bound to this connection.
func (c *pgConn) Prepare(query string) (driver.Stmt, error) {
	if c.io.IOError != nil {
		return nil, driver.ErrBadConn
	}
	return NewPgStmt(c.io, query, c.dsn)
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package maintains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
func (c *pgConn) Close() error {
	return nil
}

// Begin starts and returns a new transaction.
//
// Deprecated: Drivers should implement ConnBeginTx instead (or additionally).
func (c *pgConn) Begin() (driver.Tx, error) {
	return &PgTx{}, nil
}

// NamedValueChecker可以可选地由Conn或Stmt实现。 它为驱动程序提供了更多控制来处理Go和数据库类型，超出了允许的默认值类型。
//
//sql包按以下顺序检查值检查器，在第一个找到的匹配项处停止：
// Stmt.NamedValueChecker，Conn.NamedValueChecker，Stmt.ColumnConverter，DefaultParameterConverter。
//
//如果CheckNamedValue返回ErrRemoveArgument，则NamedValue将不包含在最终查询参数中。 这可用于将特殊选项传递给查询本身。
//
//如果返回ErrSkip，则会将列转换器错误检查路径用于参数。 司机可能希望在他们用完特殊情况后退回ErrSkip。
func (c *pgConn) CheckNamedValue(nv *driver.NamedValue) error {
	log.Println("nv.Name", nv.Name, "nv.Value", nv.Value, "nv.Ordinal", nv.Ordinal)
	log.Println(reflect.TypeOf(nv.Value))

	switch nv.Value.(type) {

	// int
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		nv.Value, _ = strconv.ParseInt(fmt.Sprintf("%d", nv.Value), 0, 64)
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *uintptr:
		nv.Value, _ = strconv.ParseInt(fmt.Sprintf("%d", reflect.ValueOf(nv.Value).Elem().Int()), 0, 64)
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []uintptr:
		var as = fmt.Sprintf("%v", nv.Value)
		nv.Value = "{" + strings.Replace(as[1:len(as)-1], " ", ",", -1) + "}"
	case *[]int, *[]int8, *[]int16, *[]int32, *[]int64, *[]uint, *[]uint8, *[]uint16, *[]uint32, *[]uint64, *[]uintptr:
		var as = fmt.Sprintf("%v", nv.Value)
		nv.Value = "{" + strings.Replace(as[2:len(as)-1], " ", ",", -1) + "}"

		// bool
	case bool:
		if nv.Value.(bool) {
			nv.Value = "t"
		} else {
			nv.Value = "f"
		}
	case *bool:
		if reflect.ValueOf(nv.Value).Elem().Bool() {
			nv.Value = "t"
		} else {
			nv.Value = "f"
		}
	case []bool:
		var as = fmt.Sprintf("%v", nv.Value)
		nv.Value = "{" + strings.Replace(as[1:len(as)-1], " ", ",", -1) + "}"
	case *[]bool:
		var as = fmt.Sprintf("%v", nv.Value)
		nv.Value = "{" + strings.Replace(as[2:len(as)-1], " ", ",", -1) + "}"

		// string
	case *string:
		nv.Value = reflect.ValueOf(nv.Value).Elem().String()
	case []string:

	}
	return nil
}

func (c *pgConn) cancel() {
	_ = c.io.CancelRequest(c.dsn.Address())
}
