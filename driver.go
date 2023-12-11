package pg

import (
	"github.com/blusewang/pg/v2/internal/app"
	"github.com/blusewang/pg/v2/internal/client/frame"
)

type Error struct {
	frame.PgError
}
type Driver struct {
	app.Driver
}
