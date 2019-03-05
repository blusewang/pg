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

type bluse struct {
	Id       int64          `json:"id"`
	Name     sql.NullString `json:"name"`
	Info     string         `json:"info"`
	CreateAt time.Time      `json:"create_at"`
	Price    float64        `json:"price"`
	UuId     string
	Raws     []byte
	IntArr   []int64
	StrArr   []string
}

func TestDriver_Open(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	db, err := sql.Open("postgres", "postgresql://bluse:@localhost/bluse?application_name=test")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	//stmt, err := db.Prepare("update bluse set arr_int=$1 where id=$2")
	stmt, err := db.Prepare("select * from bluse where id=$1")
	if err != nil {
		t.Error(err)
	}
	//rs,err := stmt.Exec([]int64{4,5,6},2)

	var a = 2
	rows, err := stmt.Query(&a)
	if err != nil {
		t.Error(err)
	}
	//log.Println(rs)
	log.Println(rows)
	//var list []bluse
	//for rows.Next() {
	//	var b bluse
	//	err = rows.Scan(&b.Id, &b.Name, &b.Info, &b.CreateAt, &b.Price, &b.UuId, &b.Raws,&b.IntArr,&b.StrArr)
	//	log.Println(b, err)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	list = append(list, b)
	//}
	//log.Println(list)
}
