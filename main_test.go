// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package pg

import (
	"bytes"
	"database/sql"
	"log"
	"testing"
	"time"
)

type bluse struct {
	Id       int64          `json:"id"`
	Name     sql.NullString `json:"name"`
	Info     string         `json:"info"`
	CreateAt time.Time      `json:"create_at"`
	Price    float64        `json:"price"`
	UuId     string
	Raws     []byte
}

func TestDriver_Open(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	log.Println(bytes.TrimLeft([]byte{0, 0, 0, 0, 0, 3, 255}, "\000"))

	db, err := sql.Open("postgres", "postgresql://bluse:@localhost/bluse?application_name=test")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	log.Println("db -> ", db)
	stmt, err := db.Prepare("select * from bluse where id<$1")
	if err != nil {
		t.Error(err)
	}
	log.Println(stmt)
	rows, err := stmt.Query(50)
	if err != nil {
		t.Error(err)
	}
	log.Println(rows)
	log.Println(rows.ColumnTypes())
	log.Println(rows.Columns())
	var list []bluse
	for rows.Next() {
		var b bluse
		err = rows.Scan(&b.Id, &b.Name, &b.Info, &b.CreateAt, &b.Price, &b.UuId, &b.Raws)
		log.Println(b, err)
		if err != nil {
			log.Println(err)
		}
		list = append(list, b)
	}
	log.Println(list)
}
