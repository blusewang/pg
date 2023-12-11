package app

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/blusewang/pg/v2/internal/client"
	"go/types"
	"strconv"
	"time"
)

type Statement struct {
	cn       *Connect
	Id       string
	SQL      string
	Response client.ParseResponse
}

func (s Statement) Close() error {
	delete(s.cn.statements, s.Id)
	return nil
}

func (s Statement) NumInput() int {
	if s.Response.Parameters == nil {
		return 0
	}
	return len(s.Response.Parameters.TypeOIDs)
}

func (s Statement) Exec(args []driver.Value) (driver.Result, error) {
	nvs := make([]driver.NamedValue, 0)
	for i, arg := range args {
		nvs = append(nvs, driver.NamedValue{
			Name:    "",
			Ordinal: i + 1,
			Value:   arg,
		})
	}
	return s.ExecContext(context.Background(), nvs)
}

func (s Statement) Query(args []driver.Value) (driver.Rows, error) {
	nvs := make([]driver.NamedValue, 0)
	for i, arg := range args {
		nvs = append(nvs, driver.NamedValue{
			Name:    "",
			Ordinal: i + 1,
			Value:   arg,
		})
	}
	return s.QueryContext(context.Background(), nvs)
}

func (s Statement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		response, err := s.cn.client.BindExec(s.Id, s.nameValue2Raw(args))
		return Result{response}, err
	}
}

func (s Statement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	select {
	case <-ctx.Done():
		return nil, sql.ErrConnDone
	default:
		response, err := s.cn.client.BindExec(s.Id, s.nameValue2Raw(args))
		return &Rows{s.cn.client.Location, s.Response.Rows, response.DataRows, 0}, err
	}
}

func (s Statement) nameValue2Raw(args []driver.NamedValue) (vs []driver.Value) {
	vs = make([]driver.Value, 0)
	for _, arg := range args {
		vs = append(vs, arg.Value)
	}
	return
}

func (s Statement) value2Row(value driver.Value) []byte {
	switch value.(type) {
	case types.Nil:
		return nil
	case int64:
		return strconv.AppendInt(nil, value.(int64), 10)
	case float64:
		return strconv.AppendFloat(nil, value.(float64), 'f', -1, 64)
	case bool:
		return strconv.AppendBool(nil, value.(bool))
	case []byte:
		return value.([]byte)
	case string:
		return []byte(value.(string))
	case time.Time:
		return []byte(value.(time.Time).Format(time.RFC3339Nano))
	case *time.Time:
		return []byte(value.(*time.Time).Format(time.RFC3339Nano))
	case json.RawMessage:
		return value.([]byte)
	default:
		return value.([]byte)
	}
}

var _ driver.Stmt = new(Statement)
var _ driver.StmtQueryContext = new(Statement)
var _ driver.StmtExecContext = new(Statement)
