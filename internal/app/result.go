package app

import (
	"database/sql/driver"
	"github.com/blusewang/pg/v2/internal/client"
)

type Result struct {
	client.BindExecResponse
}

func (r Result) LastInsertId() (int64, error) {
	return 0, nil
}

func (r Result) RowsAffected() (int64, error) {
	return int64(r.Completion.Affected()), nil
}

var _ driver.Result = new(Result)
