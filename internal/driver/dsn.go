// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"fmt"
	"net/url"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type DataSourceName struct {
	host           string
	port           string
	password       string
	connectTimeout time.Duration
	Parameter      map[string]string
	IsStrict       bool
}

func ParseDSN(connectStr string) (dsn *DataSourceName, err error) {
	dsn = new(DataSourceName)
	dsn.setDefault()
	if strings.Contains(connectStr, "://") {
		err = dsn.parseURI(connectStr)
	} else {
		err = dsn.parseDSN(connectStr)
	}
	return
}

func (dsn *DataSourceName) setDefault() {
	if runtime.GOOS == "windows" {
		dsn.host = "localhost"
	} else {
		dsn.host = "/tmp"
	}
	dsn.port = "5432"
	dsn.Parameter = make(map[string]string)
	dsn.Parameter["user"] = "postgres"
	dsn.Parameter["database"] = "postgres"
	dsn.Parameter["DateStyle"] = "ISO, YMD"
	dsn.Parameter["client_encoding"] = "UTF8"
	dsn.connectTimeout = time.Duration(60) * time.Second
	u, err := user.Current()
	if err == nil {
		dsn.Parameter["user"] = u.Name
	}
}

func (dsn *DataSourceName) parseDSN(str string) (err error) {
	p := make(map[string]string)
	for _, item := range strings.Split(str, " ") {
		pair := strings.Split(item, "=")
		if len(pair) == 2 {
			p[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	}
	if host, has := p["host"]; has {
		dsn.host = host
		delete(p, "host")
	}
	if port, has := p["port"]; has {
		dsn.port = port
		delete(p, "port")
	}
	if u, has := p["user"]; has {
		dsn.Parameter["user"] = u
		delete(p, "user")
	}
	if password, has := p["password"]; has {
		dsn.password = password
		delete(p, "password")
	}
	if dbName, has := p["dbname"]; has {
		dsn.Parameter["database"] = dbName
		delete(p, "dbname")
	}
	if appName, has := p["application_name"]; has {
		dsn.Parameter["application_name"] = appName
		delete(p, "application_name")
	} else if appName, has := p["fallback_application_name"]; has {
		dsn.Parameter["application_name"] = appName
		delete(p, "fallback_application_name")
	}
	if tos, has := p["connect_timeout"]; has {
		to, err := strconv.Atoi(tos)
		if err != nil {
			return err
		}
		dsn.connectTimeout = time.Duration(to) * time.Second
		delete(p, "connect_timeout")
	}
	if host, has := p["host"]; has {
		dsn.host = host
		delete(p, "host")
	}
	if strict, has := p["strict"]; has {
		dsn.IsStrict = strict == "true"
		delete(p, "strict")
	}
	for k, v := range p {
		dsn.Parameter[k] = v
	}
	return
}

func (dsn *DataSourceName) parseURI(uri string) (err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" && u.Scheme != "pg" {
		err = fmt.Errorf("invalid connection protocol: %s", u.Scheme)
		return
	}
	if u.Hostname() != "" {
		dsn.host = u.Hostname()
	}
	if u.Port() != "" {
		dsn.port = u.Port()
	}
	if u.User.Username() != "" {
		dsn.Parameter["user"] = u.User.Username()
	}
	if pwd, has := u.User.Password(); has {
		dsn.password = pwd
	}
	if u.Path != "" {
		ps := strings.Split(u.Path, "/")
		if len(ps) > 0 {
			dsn.Parameter["database"] = ps[1]
		}
	}
	query := u.Query()

	if tos := query.Get("connect_timeout"); tos != "" {
		to, err := strconv.Atoi(tos)
		if err != nil {
			return err
		}
		dsn.connectTimeout = time.Duration(to) * time.Second
		query.Del("connect_timeout")
	}

	for k, v := range query {
		if len(v) == 1 {
			if k == "fallback_application_name" && dsn.Parameter["application_name"] == "" {
				dsn.Parameter["application_name"] = v[0]
			} else if k == "strict" {
				dsn.IsStrict = v[0] == "true"
			} else {
				dsn.Parameter[k] = v[0]
			}
		}
	}

	return
}

func (dsn *DataSourceName) Address() (network, address string, timeout time.Duration) {
	if strings.HasPrefix(dsn.host, "/") {
		network = "unix"
		address = dsn.host + "/.s.PGSQL." + dsn.port
	} else {
		network = "tcp"
		address = dsn.host + ":" + dsn.port
	}
	timeout = dsn.connectTimeout
	return
}
