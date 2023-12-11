package pg

import (
	"github.com/blusewang/pg/internal/app"
	"github.com/blusewang/pg/internal/client/frame"
)

type Error struct {
	frame.PgError
}
type Driver struct {
	app.Driver
}
