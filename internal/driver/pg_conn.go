// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blusewang/pg/internal/client"
	"github.com/blusewang/pg/internal/helper"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type PgConn struct {
	dsn    *helper.DataSourceName
	client *client.Client
}

func NewPgConn(name string) (c PgConn, err error) {
	c.dsn, err = helper.ParseDSN(name)
	if err != nil {
		return
	}

	c.client, err = client.NewClient(context.Background(), *c.dsn)
	if err != nil {
		return
	}
	return
}

func NewPgConnContext(ctx context.Context, name string) (c *PgConn, err error) {
	c = new(PgConn)
	c.dsn, err = helper.ParseDSN(name)
	if err != nil {
		return
	}

	c.client, err = client.NewClient(ctx, *c.dsn)
	if err != nil {
		return
	}
	return
}

func (c PgConn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	if c.client.IOError != nil {
		return nil, driver.ErrBadConn
	}
	return NewPgStmtContext(ctx, c.client, query)
}

func (c PgConn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	if c.client.IOError != nil {
		return nil, driver.ErrBadConn
	}
	if c.client.IsInTransaction() {
		err = errors.New("this connection is in transaction")
	}
	_, err = c.client.QueryNoArgs("begin")
	if err != nil {
		return
	}
	if !c.client.IsInTransaction() {
		err = errors.New("begin fail")
	}

	return &PgTx{pgConn: c}, nil
}

func (c PgConn) Begin() (_ driver.Tx, err error) {
	if c.client.IOError != nil {
		return nil, driver.ErrBadConn
	}
	if c.client.IsInTransaction() {
		err = errors.New("this connection is in transaction")
	}
	_, err = c.client.QueryNoArgs("begin")
	if err != nil {
		return
	}
	if !c.client.IsInTransaction() {
		err = errors.New("begin fail")
	}

	return &PgTx{pgConn: c}, nil
}

// Prepare returns a prepared statement, bound to this connection.
func (c PgConn) Prepare(query string) (driver.Stmt, error) {
	if c.client.IOError != nil {
		return nil, driver.ErrBadConn
	}
	return NewPgStmt(c.client, query)
}

func (c PgConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (res driver.Result, err error) {
	if strings.HasPrefix(strings.ToLower(query), "listen") && len(args) == 0 {
		response, err := c.client.QueryNoArgs(query)
		res = driver.RowsAffected(response.Completion.Affected())
		return res, err
	}
	stmt, err := NewPgStmt(c.client, query)
	// 在判断 err 是否为 null前定义defer方法。
	defer func() {
		if err != nil {
			delete(c.client.StatementMaps, stmt.Id)
			_ = c.client.CloseParse(stmt.Id)
		}
	}()
	if err != nil {
		return nil, err
	}
	return stmt.ExecContext(ctx, args)
}

func (c PgConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (_ driver.Rows, err error) {
	stmt, err := NewPgStmt(c.client, query)
	// 在判断 err 是否为 null前定义defer方法。
	defer func() {
		if err != nil {
			delete(c.client.StatementMaps, stmt.Id)
			_ = c.client.CloseParse(stmt.Id)
		}
	}()
	if err != nil {
		return nil, err
	}
	return stmt.QueryContext(ctx, args)
}

// CheckNamedValue NamedValueChecker可以可选地由Conn或Stmt实现。 它为驱动程序提供了更多控制来处理Go和数据库类型，超出了允许的默认值类型。
//
//sql包按以下顺序检查值检查器，在第一个找到的匹配项处停止：
// Stmt.NamedValueChecker，Conn.NamedValueChecker，Stmt.ColumnConverter，DefaultParameterConverter。
//
//如果CheckNamedValue返回ErrRemoveArgument，则NamedValue将不包含在最终查询参数中。 这可用于将特殊选项传递给查询本身。
//
//如果返回ErrSkip，则会将列转换器错误检查路径用于参数。 司机可能希望在他们用完特殊情况后退回ErrSkip。
func (c PgConn) CheckNamedValue(nv *driver.NamedValue) error {
	if nv.Value == nil {
		return nil
	} else {
		var vi = reflect.ValueOf(nv.Value)
		if vi.Kind() == reflect.Ptr && vi.IsNil() {
			nv.Value = nil
			return nil
		}
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
	case string:
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
		raw, _ := json.Marshal(nv.Value)
		var as = string(raw)
		as = strings.ReplaceAll(as, "[", "{")
		nv.Value = strings.ReplaceAll(as, "]", "}")
	case *[]float32, *[]float64:
		raw, _ := json.Marshal(nv.Value)
		var as = string(raw)
		as = strings.ReplaceAll(as, "[", "{")
		nv.Value = strings.ReplaceAll(as, "]", "}")
	case [][]float32, [][]float64:
		raw, _ := json.Marshal(nv.Value)
		var as = string(raw)
		as = strings.ReplaceAll(as, "[", "{")
		nv.Value = strings.ReplaceAll(as, "]", "}")

	//	byte
	case []byte:
		nv.Value = fmt.Sprintf("\\x%x", nv.Value)
	case *[]byte:
		var as = fmt.Sprintf("\\x%x", nv.Value)
		nv.Value = "\\x" + as[1:]

	// time
	case time.Time:
	case *time.Time:

	default:
		nv.Value = fmt.Sprintf("%v", nv.Value)
	}
	return nil
}

func (c PgConn) cancel() {
	_ = c.client.CancelRequest()
}

func (c PgConn) Close() (err error) {
	err = c.client.Terminate()
	return
}
