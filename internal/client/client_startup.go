// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/blusewang/pg/internal/frame"
	"github.com/blusewang/pg/internal/helper"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

func (c *Client) connect(ctx context.Context, dsn helper.DataSourceName) (err error) {
	nw, addr, timeout := dsn.Address()
	d := net.Dialer{Timeout: timeout}
	c.conn, err = d.DialContext(ctx, nw, addr)
	if err != nil {
		log.Println(ctx, nw, addr, err)
		return
	}
	if dsn.SSL.Mode != "disable" && dsn.SSL.Mode != "allow" {
		if err = c.ssl(dsn); err != nil {
			return
		}
	}
	c.writer = frame.NewEncoder(c.conn)
	c.reader = frame.NewDecoder(c.conn)

	return c.startup(dsn)
}

func (c *Client) ssl(dsn helper.DataSourceName) (err error) {
	var tlsConfig tls.Config
	if err = frame.NewEncoder(c.conn).Fire(frame.NewSSLRequest()); err != nil {
		return
	}
	var code byte
	code, err = bufio.NewReader(c.conn).ReadByte()
	if err != nil {
		log.Println(code, err)
		return
	}

	switch dsn.SSL.Mode {
	case "prefer":
		if code == 'N' {
			// 服务器不支持，转用明文传输
			return nil
		} else if err = dsn.SSLCheck(); err != nil {
			// 服务器支持，但本地证书检查不通过，也转用明文传输
			return nil
		}
		tlsConfig.InsecureSkipVerify = true
	case "require":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		// 如果没有根证书，标记为只验证证书
		if _, err = os.Stat(dsn.SSL.RootCert); err == nil {
			tlsConfig.InsecureSkipVerify = true
		}
	case "verify-ca":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		tlsConfig.InsecureSkipVerify = true
	case "verify-full":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		tlsConfig.ServerName = dsn.Host
	default:
		return errors.New("pq: SSL unknown type")
	}

	if err = dsn.SSLCheck(); err != nil {
		return
	}
	tlsConfig.Renegotiation = tls.RenegotiateFreelyAsClient
	cert, err := tls.LoadX509KeyPair(dsn.SSL.Cert, dsn.SSL.Key)
	if err != nil {
		return err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	if _, err = os.Stat(dsn.SSL.RootCert); err == nil {
		tlsConfig.RootCAs = x509.NewCertPool()
		var raw []byte
		raw, err = ioutil.ReadFile(dsn.SSL.RootCert)
		if err != nil {
			return err
		}
		if !tlsConfig.RootCAs.AppendCertsFromPEM(raw) {
			return errors.New("pg: can't parse root cert")
		}
	}

	// 升级至TLS
	c.conn = tls.Client(c.conn, &tlsConfig)
	c.writer = frame.NewEncoder(c.conn)
	c.reader = frame.NewDecoder(c.conn)
	return
}

func (c *Client) startup(dsn helper.DataSourceName) (err error) {
	su := frame.NewStartup()
	for k, v := range dsn.Parameter {
		su.AddParam(k, v)
	}
	su.WriteUint8(0)
	if err = c.writer.Fire(su.Frame); err != nil {
		return
	}
	var out interface{}
	for {
		out, err = c.reader.Decode()
		if err != nil {
			return
		}
		switch f := out.(type) {
		case *frame.AuthRequest:
			switch f.GetType() {
			case frame.AuthTypePwd:
				ar := frame.NewAuthResponse()
				ar.Password(dsn.Password)
				if err = c.writer.Fire(ar.Frame); err != nil {
					return
				}
			case frame.AuthTypeMd5:
				ar := frame.NewAuthResponse()
				ar.Md5Pwd(dsn.Parameter["user"], dsn.Password, string(f.GetMd5Salt()))
				if err = c.writer.Fire(ar.Frame); err != nil {
					return
				}
			case frame.AuthTypeOk:
			}
		case *frame.ParameterStatus:
			f.Decode()
			c.parameterStatus[f.Name] = f.Value
			if f.Name == "TimeZone" {
				c.Location, err = time.LoadLocation(f.Value)
				if err != nil {
					c.Location = nil
				}
			}
		case *frame.BackendKeyData:
			f.Decode()
			c.backendPid = f.Pid
			c.backendKey = f.Key
		case *frame.ReadyForQuery:
			c.status = frame.TransactionStatus(f.Payload[0])
			go c.readerLoop()
			return nil
		}
	}
}

func (c *Client) readerLoop() {
	var out interface{}
	var err error
	for {
		out, err = c.reader.Decode()
		if err != nil {
			_ = c.Terminate()
			return
		}
		switch f := out.(type) {
		case *frame.Notification:
			f.Decode()
			if ListenMap[f.Condition] != nil {
				ListenMap[f.Condition](f.Text)
			}
		case *frame.NoticeResponse:
			f.Decode()
			raw, _ := json.Marshal(f.Error.Error)
			log.Println(string(raw))
		case *frame.Error:
			f.Decode()
			c.Err = f
			c.status = frame.TransactionStatusNoReady
			if f.Error.Fail == "FATAL" || f.Error.Fail == "PANIC" {
				_ = c.Terminate()
				return
			}
		case *frame.ReadyForQuery:
			c.status = frame.TransactionStatus(f.Payload[0])
			c.frameChan <- f
		default:
			c.frameChan <- f
		}
	}
}
