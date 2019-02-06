// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package pg

import (
	"database/sql"
	"log"
	"testing"
	"time"
)

func TestDriver_Open(t *testing.T) {
	l, err := time.LoadLocation("PRC")
	log.Println(l, err)

	db, err := sql.Open("postgres", "postgresql://bluse:@localhost/bluse?application_name=test")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	log.Println("db -> ", db)
	rows, err := db.Query("select * from bluse")
	if err != nil {
		t.Error(err)
	}
	log.Println("rows -> ", rows)
}
