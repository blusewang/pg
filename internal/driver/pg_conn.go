// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/blusewang/pg/internal/network"
	"reflect"
	"strconv"
	"strings"
)

func NewPgConn(name string) (c *PgConn, err error) {
	c = new(PgConn)
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
	c.stmts = make(map[string]*PgStmt)
	return
}

func NewPgConnContext(ctx context.Context, name string) (c *PgConn, err error) {
	c = new(PgConn)
	c.dsn, err = ParseDSN(name)
	if err != nil {
		return
	}

	c.io = network.NewPgIO()
	var net, addr, timeout = c.dsn.Address()
	err = c.io.DialContext(ctx, net, addr, timeout)
	if err != nil {
		return
	}
	err = c.io.StartUp(c.dsn.Parameter, c.dsn.password)
	if err != nil {
		return
	}
	c.stmts = make(map[string]*PgStmt)
	return
}

type PgConn struct {
	dsn   *DataSourceName
	io    *network.PgIO
	stmts map[string]*PgStmt
}

// Prepare returns a prepared statement, bound to this connection.
func (c *PgConn) Prepare(query string) (driver.Stmt, error) {
	if c.io.IOError != nil {
		return nil, driver.ErrBadConn
	}
	return NewPgStmt(c, query)
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package maintains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
func (c *PgConn) Close() (err error) {
	err = c.io.Terminate()
	return
}

// Begin starts and returns a new transaction.
//
// Deprecated: Drivers should implement ConnBeginTx instead (or additionally).
func (c *PgConn) Begin() (_ driver.Tx, err error) {
	if c.io.IOError != nil {
		return nil, driver.ErrBadConn
	}
	if c.io.IsInTransaction() {
		err = errors.New("this connection is in transaction")
	}
	_, _, _, err = c.io.QueryNoArgs("begin")
	if err != nil {
		return
	}
	if !c.io.IsInTransaction() {
		err = errors.New("begin fail")
	}

	return &PgTx{pgConn: c}, nil
}

func (c *PgConn) Query(query string, args []driver.Value) (_ driver.Rows, err error) {
	stmt, err := NewPgStmt(c, query)
	if err != nil {
		return nil, err
	}
	return stmt.Query(args)
}

// NamedValueChecker可以可选地由Conn或Stmt实现。 它为驱动程序提供了更多控制来处理Go和数据库类型，超出了允许的默认值类型。
//
//sql包按以下顺序检查值检查器，在第一个找到的匹配项处停止：
// Stmt.NamedValueChecker，Conn.NamedValueChecker，Stmt.ColumnConverter，DefaultParameterConverter。
//
//如果CheckNamedValue返回ErrRemoveArgument，则NamedValue将不包含在最终查询参数中。 这可用于将特殊选项传递给查询本身。
//
//如果返回ErrSkip，则会将列转换器错误检查路径用于参数。 司机可能希望在他们用完特殊情况后退回ErrSkip。
func (c *PgConn) CheckNamedValue(nv *driver.NamedValue) error {
	if nv.Value == nil {
		return nil
	} else if !reflect.ValueOf(nv.Value).Elem().CanAddr() {
		nv.Value = nil
		return nil
	}
	switch nv.Value.(type) {

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

	//	string
	case *string:
		nv.Value = reflect.ValueOf(nv.Value).Elem().String()
	case []string:
		if len(nv.Value.([]string)) > 0 {
			var str = "{"
			for _, s := range nv.Value.([]string) {
				s = strings.Replace(s, `"`, `\"`, -1)
				s = strings.Replace(s, `'`, `\'`, -1)
				if strings.Contains(s, ",") {
					str += `"` + s + `",`
				} else {
					str += s + ","
				}
			}
			nv.Value = str[:len(str)-1] + "}"
		} else {
			nv.Value = "{}"
		}
	case *[]string:
		if nv.Value == nil {
			return nil
		}
		if len(*nv.Value.(*[]string)) > 0 {
			var str = "{"
			for _, s := range *nv.Value.(*[]string) {
				s = strings.Replace(s, `"`, `\"`, -1)
				s = strings.Replace(s, `'`, `\'`, -1)
				if strings.Contains(s, ",") {
					str += `"` + s + `",`
				} else {
					str += s + ","
				}
			}
			nv.Value = str[:len(str)-1] + "}"
		} else {
			nv.Value = "{}"
		}

	//	int
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint64:
		nv.Value, _ = strconv.ParseInt(fmt.Sprintf("%d", nv.Value), 0, 64)
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint64:
		nv.Value, _ = strconv.ParseInt(fmt.Sprintf("%d", reflect.ValueOf(nv.Value).Elem().Int()), 0, 64)
	case []int, []int8, []int16, []int64, []uint, []uint16, []uint64:
		var as = fmt.Sprintf("%v", nv.Value)
		nv.Value = "{" + strings.Replace(as[1:len(as)-1], " ", ",", -1) + "}"
	case *[]int, *[]int8, *[]int16, *[]int64, *[]uint, *[]uint16, *[]uint64:
		var as = fmt.Sprintf("%v", nv.Value)
		nv.Value = "{" + strings.Replace(as[2:len(as)-1], " ", ",", -1) + "}"

	//	float
	case float32, float64:
		nv.Value, _ = strconv.ParseFloat(fmt.Sprintf("%v", nv.Value), 64)
	case *float32, *float64:
		nv.Value, _ = strconv.ParseFloat(fmt.Sprintf("%v", reflect.ValueOf(nv.Value).Elem().Float()), 64)
	case []float32, []float64:
		var as = fmt.Sprintf("%v", nv.Value)
		nv.Value = "{" + strings.Replace(as[1:len(as)-1], " ", ",", -1) + "}"
	case *[]float32, *[]float64:
		var as = fmt.Sprintf("%v", nv.Value)
		nv.Value = "{" + strings.Replace(as[2:len(as)-1], " ", ",", -1) + "}"

	//	byte
	case []byte:
		nv.Value = fmt.Sprintf("\\x%x", nv.Value)
	case *[]byte:
		var as = fmt.Sprintf("\\x%x", nv.Value)
		nv.Value = "\\x" + as[1:]

	default:
		return driver.ErrRemoveArgument

	}
	return nil
}

func (c *PgConn) cancel() {
	_ = c.io.CancelRequest(c.dsn.Address())
}
