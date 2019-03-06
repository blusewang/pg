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

func TestDriver_Query(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	db, err := sql.Open("pg", "postgresql://bluse:@localhost/bluse?application_name=test")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from bluse where id>$1", 0)
	if err != nil {
		t.Error(err)
	}

	for rows.Next() {
		var b bluse
		err = rows.Scan(&b.Id, &b.Name, &b.Info, &b.CreateAt, &b.Price, &b.UuId, &b.Raws, &b.IntArr, &b.StrArr)
		log.Println(b, err)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(b.Raws))
	}
}

func TestDriver_TransactionQuery(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	db, err := sql.Open("pg", "postgresql://bluse:@localhost/bluse?application_name=test")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}
	rows, err := tx.Query("select * from bluse where id>$1", 0)
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		var b bluse
		err = rows.Scan(&b.Id, &b.Name, &b.Info, &b.CreateAt, &b.Price, &b.UuId, &b.Raws, &b.IntArr, &b.StrArr)
		log.Println(b, err)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(b.Raws))
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}

func TestDriver_TransactionExec(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	db, err := sql.Open("pg", "postgresql://bluse:@localhost/bluse?application_name=test")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}
	rs, err := tx.Exec("update bluse set price=$1 where id=$2", 4.13, 1)
	if err != nil {
		t.Error(err)
	}
	log.Println(rs.RowsAffected())
	log.Println(rs.LastInsertId())
	err = tx.Rollback()
	if err != nil {
		t.Error(err)
	}
}

func TestDriver_Select(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	db, err := sql.Open("pg", "postgresql://bluse:@localhost/bluse?application_name=test")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from bluse where id>$1")
	if err != nil {
		t.Error(err)
	}
	stmt2, err := db.Prepare("select * from bluse where id>$1")
	log.Println(stmt2, err)
	if err != nil {
		t.Error(err)
	}

	var a = 0
	rows, err := stmt.Query(&a)
	if err != nil {
		t.Error(err)
	}
	//log.Println(rows)
	//var list []bluse
	for rows.Next() {
		var b bluse
		err = rows.Scan(&b.Id, &b.Name, &b.Info, &b.CreateAt, &b.Price, &b.UuId, &b.Raws, &b.IntArr, &b.StrArr)
		log.Println(b, err)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(b.Raws))
	}
	//log.Println(list)
}

func TestDriver_Exec(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	db, err := sql.Open("pg", "postgresql://bluse:@localhost/bluse?application_name=test")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update bluse set raws=$1 where id=$2")
	if err != nil {
		t.Error(err)
	}
	rs, err := stmt.Exec([]byte("asdf,'sdf"), 1)

	log.Println(rs.LastInsertId())
	log.Println(rs.RowsAffected())
}

func TestDriver_OpenDb(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	var db = sql.OpenDB(NewConnector(""))
	log.Println(db)
}
