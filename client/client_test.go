// Copyright 2022 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"context"
	"database/sql/driver"
	"github.com/blusewang/pg/dsn"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	d, err := dsn.ParseDSN("pg://developer:dev.123@r/core")
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewClient(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	s, err := c.Parse("abc", "select * from pg_settings where name=$1")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(s)
	res, err := c.ParseExec("abc", []driver.Value{"allow_system_table_mods"})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res)
	log.Println(c.Close())
}
