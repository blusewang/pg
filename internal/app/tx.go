package app

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/blusewang/pg/internal/client"
)

type Tx struct {
	client *client.Client
	ctx    context.Context
}

func (t Tx) Commit() (err error) {
	select {
	case <-t.ctx.Done():
		return
	default:
		if t.client.IsInTransaction() == false {
			err = errors.New("this connection is out of transaction")
			return
		}
		_, err = t.client.QueryNoArgs("commit")
		if err != nil {
			return
		}
		if t.client.IsInTransaction() {
			err = errors.New("commit fail")
		}
		return
	}
}

func (t Tx) Rollback() (err error) {
	select {
	case <-t.ctx.Done():
		return
	default:
		if t.client.IsInTransaction() == false {
			err = errors.New("this connection is out of transaction")
			return
		}
		_, err = t.client.QueryNoArgs("rollback")
		if err != nil {
			return
		}
		if t.client.IsInTransaction() {
			err = errors.New("rollback fail")
		}
		return
	}
}

var _ driver.Tx = new(Tx)
