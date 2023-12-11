package pg

import (
	"context"
	"github.com/blusewang/pg/v2/internal/client"
)

type Listener interface {
	Listen(channel string) error
	GetNotification() (pid uint32, channel, message string, err error)
	Terminate() error
}

func NewListener(ctx context.Context, dsnString string) (Listener, error) {
	c := client.NewClient()
	dsn, err := client.ParseDSN(dsnString)
	if err != nil {
		return nil, err
	}
	if err = c.Connect(ctx, dsn); err != nil {
		return nil, err
	}
	if err = c.AutoSSL(); err != nil {
		return nil, err
	}
	if err = c.Startup(); err != nil {
		return nil, err
	}
	return c, nil
}
