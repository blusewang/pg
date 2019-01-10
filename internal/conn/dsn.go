// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package conn

import (
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type dataSourceName struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	ApplicationName string
	ConnectTimeout  time.Duration
}

func ParseDSN(connectStr string) (dsn dataSourceName, err error) {
	if strings.Contains(connectStr, "://") {
		err = dsn.parseURI(connectStr)
	} else {
		err = dsn.parseDSN(connectStr)
	}
	return
}

func (dsn *dataSourceName) setDefault() {
	if runtime.GOOS == "windows" {
		dsn.Host = "localhost"
	} else {
		dsn.Host = "/tmp"
	}
	dsn.Port = "5432"
	dsn.User = "postgres"
	dsn.DbName = "postgres"
	dsn.ConnectTimeout = time.Duration(60) * time.Second
}

func (dsn *dataSourceName) parseDSN(str string) (err error) {
	p := make(map[string]string)
	for _, item := range strings.Split(str, " ") {
		pair := strings.Split(item, "=")
		if len(pair) == 2 {
			p[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	}
	dsn.setDefault()
	if host, has := p["host"]; has {
		dsn.Host = host
	}
	if port, has := p["port"]; has {
		dsn.Port = port
	}
	if user, has := p["user"]; has {
		dsn.User = user
	}
	if password, has := p["password"]; has {
		dsn.Password = password
	}
	if dbName, has := p["dbname"]; has {
		dsn.DbName = dbName
	}
	if applicationName, has := p["application_name"]; has {
		dsn.ApplicationName = applicationName
	}
	if tos, has := p["connect_timeout"]; has {
		to, err := strconv.Atoi(tos)
		if err != nil {
			return err
		}
		dsn.ConnectTimeout = time.Duration(to) * time.Second
	}
	if host, has := p["host"]; has {
		dsn.Host = host
	}
	return
}

func (dsn *dataSourceName) parseURI(uri string) (err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		err = fmt.Errorf("invalid connection protocol: %s", u.Scheme)
		return
	}
	dsn.setDefault()
	if u.Hostname() != "" {
		dsn.Host = u.Hostname()
	}
	if u.User.Username() != "" {
		dsn.User = u.User.Username()
	}
	if pwd, has := u.User.Password(); has {
		dsn.Password = pwd
	}
	if u.Path != "" {
		ps := strings.Split(u.Path, "/")
		if len(ps) > 0 {
			dsn.DbName = ps[1]
		}
	}
	dsn.ApplicationName = u.Query().Get("application_name")
	if tos := u.Query().Get("connect_timeout"); tos != "" {
		to, err := strconv.Atoi(tos)
		if err != nil {
			return err
		}
		dsn.ConnectTimeout = time.Duration(to) * time.Second
	}

	return
}

func (dsn *dataSourceName) Address() (network, address string, timeout time.Duration) {
	if strings.HasPrefix(dsn.Host, "/") {
		network = "unix"
		address = dsn.Host + ".s.PGSQL." + dsn.Port
	} else {
		network = "tcp"
		address = dsn.Host + ":" + dsn.Port
	}
	timeout = dsn.ConnectTimeout
	return
}
