// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package helper

import (
	"github.com/blusewang/pg/dsn"
	"log"
	"testing"
	"time"
)

func TestParseDSNUseURI(t *testing.T) {
	dsn, err := dsn.ParseDSN("pg://cashier:@:33521/cashier?host=/tmp&application_name=cashier_production&strict=true")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(dsn)
}

func TestParseDSNUseURIHasSock(t *testing.T) {
	name := dsn.DataSourceName{
		Host:           "/tmp",
		Port:           "5432",
		Password:       "pass.word",
		ConnectTimeout: time.Duration(10) * time.Second,
	}
	dsn, err := dsn.ParseDSN("pg://postgres:pass.word@:5432/db_name?application_name=application_name&Host=/tmp&connect_timeout=10")
	if err != nil {
		t.Fatal(err)
	}
	if dsn.Host != name.Host {
		t.Fail()
	}
	log.Println(dsn)
}

func TestParseDSNUseStr(t *testing.T) {
	name := dsn.DataSourceName{
		Host:           "postgresql.com",
		Port:           "5432",
		Password:       "pass.word",
		ConnectTimeout: time.Duration(10) * time.Second,
	}
	dsn, err := dsn.ParseDSN("user=postgres Password=pass.word Host=postgresql.com Port=5432 dbname=db_name application_name=application_name connect_timeout=10")
	if err != nil {
		t.Fatal(err)
	}
	if dsn.Host != name.Host {
		t.Fail()
	}
}

func TestParseDSNUseStrHasSock(t *testing.T) {
	name := dsn.DataSourceName{
		Host:           "/tmp",
		Port:           "5432",
		Password:       "pass.word",
		ConnectTimeout: time.Duration(10) * time.Second,
	}
	dsn, err := dsn.ParseDSN("user=postgres Password=pass.word Host=/tmp Port=5432 dbname=db_name application_name=application_name connect_timeout=10")
	if err != nil {
		t.Fatal(err)
	}
	if dsn.Host != name.Host {
		t.Fail()
	}
	log.Println(dsn)
}
