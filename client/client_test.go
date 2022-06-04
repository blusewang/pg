// Copyright 2022 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"context"
	"github.com/blusewang/pg/dsn"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	d, err := dsn.ParseDSN("pg://postgres:abc123@r:5433/core?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewClient(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.QueryNoArgs("select now()")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res)
	log.Println(c.Close())
}
