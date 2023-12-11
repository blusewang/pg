package app

import (
	"context"
	"crypto/md5"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blusewang/pg/v2/internal/client"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewConnect(ctx context.Context, dsn client.DataSourceName) (c Connect, err error) {
	c.client = client.NewClient()
	if err != nil {
		return
	}
	if err = c.client.Connect(ctx, dsn); err != nil {
		return
	}
	if err = c.client.AutoSSL(); err != nil {
		return
	}
	if err = c.client.Startup(); err != nil {
		return
	}
	c.statements = make(map[string]*Statement)
	return
}

type Connect struct {
	client     *client.Client
	statements map[string]*Statement
}

func (c Connect) IsValid() bool {
	return c.client.ConnectStatus != client.ConnectStatusDisconnected
}

func (c Connect) ResetSession(_ context.Context) (err error) {
	return nil
}

func (c Connect) query2Id(query string) string {
	raw := md5.Sum([]byte(query))
	return hex.EncodeToString(raw[:])
}

func (c Connect) parse(ctx context.Context, query string) (stmt *Statement, err error) {
	select {
	case <-ctx.Done():
		return
	default:
		id := c.query2Id(query)
		if c.statements[id] == nil {
			res, err := c.client.Parse(id, query)
			if err != nil {
				return nil, err
			}
			c.statements[id] = &Statement{cn: &c, Id: id, SQL: query, Response: res}
		}
		return c.statements[id], nil
	}
}

func (c Connect) Prepare(query string) (stmt driver.Stmt, err error) {
	return c.parse(context.Background(), query)
}

func (c Connect) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			return c.parse(ctx, query)
		}
	}
}

func (c Connect) Close() error {
	return c.client.Terminate()
}

func (c Connect) Begin() (tx driver.Tx, err error) {
	if c.client.IsInTransaction() {
		err = errors.New("this connection is in transaction")
		return
	}
	_, err = c.client.QueryNoArgs("begin")
	if err != nil {
		return
	}
	if !c.client.IsInTransaction() {
		err = errors.New("begin fail")
	}

	return &Tx{client: c.client, ctx: context.Background()}, nil
}

func (c Connect) BeginTx(ctx context.Context, _ driver.TxOptions) (tx driver.Tx, err error) {
	if c.client.IsInTransaction() {
		err = errors.New("this connection is in transaction")
		return
	}
	_, err = c.client.QueryNoArgs("begin")
	if err != nil {
		return
	}
	if !c.client.IsInTransaction() {
		err = errors.New("begin fail")
	}

	return &Tx{client: c.client, ctx: ctx}, nil
}

func (c Connect) CheckNamedValue(nv *driver.NamedValue) error {
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

	case json.RawMessage:
		nv.Value = string(nv.Value.(json.RawMessage))

	// time
	case time.Time:
	case *time.Time:

	default:
		nv.Value = fmt.Sprintf("%v", nv.Value)
	}
	return nil
}

func (c Connect) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	stmt, err := c.parse(ctx, query)
	if err != nil {
		return
	}
	return stmt.QueryContext(ctx, args)
}

func (c Connect) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (res driver.Result, err error) {
	stmt, err := c.parse(ctx, query)
	if err != nil {
		return
	}
	return stmt.ExecContext(ctx, args)
}

var _ driver.Conn = new(Connect)
var _ driver.ConnPrepareContext = new(Connect)
var _ driver.ConnBeginTx = new(Connect)
var _ driver.Validator = new(Connect)
var _ driver.NamedValueChecker = new(Connect)
var _ driver.QueryerContext = new(Connect)
var _ driver.ExecerContext = new(Connect)
var _ driver.SessionResetter = new(Connect)
