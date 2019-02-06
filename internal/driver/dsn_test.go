// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"log"
	"testing"
	"time"
)

func TestParseDSNUseURI(t *testing.T) {
	name := DataSourceName{
		host:           "postgresql.com",
		port:           "5432",
		password:       "pass.word",
		connectTimeout: time.Duration(10) * time.Second,
	}
	dsn, err := ParseDSN("postgresql://postgres:pass.word@postgresql.com:5432/db_name?application_name=application_name&connect_timeout=10")
	if err != nil {
		t.Fatal(err)
	}
	if dsn.host != name.host {
		t.Fail()
	}
	log.Println(dsn.Address())
	log.Println(dsn.Parameter)
}

func TestParseDSNUseURIHasSock(t *testing.T) {
	name := DataSourceName{
		host:           "/tmp",
		port:           "5432",
		password:       "pass.word",
		connectTimeout: time.Duration(10) * time.Second,
	}
	dsn, err := ParseDSN("postgresql://postgres:pass.word@:5432/db_name?application_name=application_name&host=/tmp&connect_timeout=10")
	if err != nil {
		t.Fatal(err)
	}
	if dsn.host != name.host {
		t.Fail()
	}
}

func TestParseDSNUseStr(t *testing.T) {
	name := DataSourceName{
		host:           "postgresql.com",
		port:           "5432",
		password:       "pass.word",
		connectTimeout: time.Duration(10) * time.Second,
	}
	dsn, err := ParseDSN("user=postgres password=pass.word host=postgresql.com port=5432 dbname=db_name application_name=application_name connect_timeout=10")
	if err != nil {
		t.Fatal(err)
	}
	if dsn.host != name.host {
		t.Fail()
	}
}

func TestParseDSNUseStrHasSock(t *testing.T) {
	name := DataSourceName{
		host:           "/tmp",
		port:           "5432",
		password:       "pass.word",
		connectTimeout: time.Duration(10) * time.Second,
	}
	dsn, err := ParseDSN("user=postgres password=pass.word host=/tmp port=5432 dbname=db_name application_name=application_name connect_timeout=10")
	if err != nil {
		t.Fatal(err)
	}
	if dsn.host != name.host {
		t.Fail()
	}
	log.Println(dsn)
}
