// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package helper

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type DataSourceName struct {
	Host           string
	Port           string
	Password       string
	ConnectTimeout time.Duration
	Parameter      map[string]string
	SSL            struct {
		Mode        string
		Cert        string
		Key         string
		RootCert    string
		Crl         string
		Compression int
	}
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
	dsn.Port = "5432"
	dsn.Parameter = make(map[string]string)
	dsn.Parameter["user"] = "postgres"
	dsn.Parameter["database"] = "postgres"
	dsn.Parameter["DateStyle"] = "ISO, YMD"
	dsn.Parameter["client_encoding"] = "UTF8"
	dsn.ConnectTimeout = time.Duration(60) * time.Second
	dsn.SSL.Compression = 1
	dsn.SSL.Mode = "prefer"
	u, err := user.Current()
	if err == nil {
		dsn.Parameter["user"] = u.Name
	}
	if runtime.GOOS == "windows" {
		dsn.Host = "localhost"
		var envs = os.Environ()
		var baseDir = dsn.pickEnv("APPDATA", envs)
		if baseDir != "" {
			dsn.SSL.Cert = dsn.pickEnv("PGSSLCERT", envs)
			if dsn.SSL.Cert == "" {
				dsn.SSL.Cert = baseDir + "\\postgresql\\postgresql.crt"
			}
			dsn.SSL.Key = dsn.pickEnv("PGSSLKEY", envs)
			if dsn.SSL.Key == "" {
				dsn.SSL.Key = baseDir + "\\postgresql\\postgresql.key"
			}
			dsn.SSL.RootCert = baseDir + "\\postgresql\\root.crt"
			dsn.SSL.Crl = baseDir + "\\postgresql\\root.crl"
		}
	} else {
		dsn.Host = "/tmp"
		if u != nil {
			dsn.SSL.Cert = filepath.Join(u.HomeDir, ".postgresql", "postgresql.crt")
			dsn.SSL.Key = filepath.Join(u.HomeDir, ".postgresql", "postgresql.key")
			dsn.SSL.RootCert = filepath.Join(u.HomeDir, ".postgresql", "root.crt")
			dsn.SSL.Crl = filepath.Join(u.HomeDir, ".postgresql", "root.crl")
		}
	}
}

func (dsn *DataSourceName) pickEnv(key string, env []string) string {
	for _, v := range env {
		if strings.HasPrefix(key, v) {
			if len(strings.Split(v, "=")) == 2 {
				return strings.Split(v, "=")[1]
			}
		}
	}
	return ""
}
func (dsn *DataSourceName) parseDSN(str string) (err error) {
	p := make(map[string]string)
	for _, item := range strings.Split(str, " ") {
		pair := strings.Split(item, "=")
		if len(pair) == 2 {
			p[strings.ToLower(strings.TrimSpace(pair[0]))] = strings.TrimSpace(pair[1])
		}
	}
	if host, has := p["host"]; has {
		dsn.Host = host
		delete(p, "host")
	}
	if port, has := p["Port"]; has {
		dsn.Port = port
		delete(p, "Port")
	}
	if u, has := p["user"]; has {
		dsn.Parameter["user"] = u
		delete(p, "user")
	}
	if password, has := p["password"]; has {
		dsn.Password = password
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
		dsn.ConnectTimeout = time.Duration(to) * time.Second
		delete(p, "connect_timeout")
	}

	dsn.pickSSLSetting(&p)

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
		dsn.Host = u.Hostname()
	}
	if u.Port() != "" {
		dsn.Port = u.Port()
	}
	if u.User.Username() != "" {
		dsn.Parameter["user"] = u.User.Username()
	}
	if pwd, has := u.User.Password(); has {
		dsn.Password = pwd
	}
	if u.Path != "" {
		ps := strings.Split(u.Path, "/")
		if len(ps) > 0 {
			dsn.Parameter["database"] = ps[1]
		}
	}
	var qm = make(map[string]string)
	for k, v := range u.Query() {
		qm[strings.ToLower(k)] = v[0]
	}

	if tos, has := qm["connect_timeout"]; has {
		to, err := strconv.Atoi(tos)
		if err != nil {
			return err
		}
		dsn.ConnectTimeout = time.Duration(to) * time.Second
		delete(qm, "connect_timeout")
	}

	if appName, has := qm["application_name"]; has {
		dsn.Parameter["application_name"] = appName
		delete(qm, "application_name")
	} else if appName, has := qm["fallback_application_name"]; has {
		dsn.Parameter["application_name"] = appName
		delete(qm, "fallback_application_name")
	}

	if host, has := qm["host"]; has {
		dsn.Host = host
		delete(qm, "host")
	}

	dsn.pickSSLSetting(&qm)

	for k, v := range qm {
		dsn.Parameter[k] = v
	}
	return
}

func (dsn *DataSourceName) pickSSLSetting(envs *map[string]string) {
	if envs != nil {
		if strings.HasPrefix(dsn.Host, "/") {
			dsn.SSL.Mode = "disable"
		} else if v, has := (*envs)["sslmode"]; has {
			dsn.SSL.Mode = v
			delete(*envs, "sslmode")
		}
		if v, has := (*envs)["sslcompression"]; has {
			if v == "0" {
				dsn.SSL.Compression = 0
			}
			delete(*envs, "sslcompression")
		}
		if v, has := (*envs)["sslcert"]; has {
			dsn.SSL.Cert = v
			delete(*envs, "sslcert")
		}
		if v, has := (*envs)["sslkey"]; has {
			dsn.SSL.Key = v
			delete(*envs, "sslkey")
		}
		if v, has := (*envs)["sslrootcert"]; has {
			dsn.SSL.RootCert = v
			delete(*envs, "sslrootcert")
		}
		if v, has := (*envs)["sslcrl"]; has {
			dsn.SSL.Crl = v
			delete(*envs, "sslcrl")
		}
	}
}

func (dsn *DataSourceName) Address() (network, address string, timeout time.Duration) {
	if strings.HasPrefix(dsn.Host, "/") {
		network = "unix"
		address = dsn.Host + "/.s.PGSQL." + dsn.Port
	} else {
		network = "tcp"
		address = dsn.Host + ":" + dsn.Port
	}
	timeout = dsn.ConnectTimeout
	return
}

func (dsn *DataSourceName) SSLCheck() (err error) {
	if dsn.SSL.Mode == "verify-ca" || dsn.SSL.Mode == "verify-full" {
		if _, err = os.Stat(dsn.SSL.RootCert); err != nil {
			return err
		}
	}

	if _, err = os.Stat(dsn.SSL.Cert); err != nil {
		return err
	}

	info, err := os.Stat(dsn.SSL.Key)
	if err != nil {
		return err
	}
	if info.Mode().Perm()&0077 != 0 {
		return errors.New("pq: Private key file has group or world access. Permissions should be u=rw (0600) or less")
	}
	return
}
