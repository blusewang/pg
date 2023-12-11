// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"bufio"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/blusewang/pg/internal/frame"
	scram2 "github.com/blusewang/pg/scram"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

func (c *Client) connect(ctx context.Context) (err error) {
	nw, addr, timeout := c.dsn.Address()
	d := net.Dialer{Timeout: timeout}
	c.conn, err = d.DialContext(ctx, nw, addr)
	if err != nil {
		log.Println(nw, addr, err)
		return
	}
	if c.dsn.SSL.Mode != "disable" && c.dsn.SSL.Mode != "allow" {
		if err = c.ssl(); err != nil {
			return
		}
	}
	c.writer = frame.NewEncoder(c.conn)
	c.reader = frame.NewDecoder(c.conn)
	c.Err = nil
	c.IOError = nil
	c.status = frame.TransactionStatusNoReady
	c.parameterMaps = make(map[string]string)
	c.StatementMaps = make(map[string]*Statement)
	c.frameChan = make(chan frame.Frame)
	c.NotifyChan = make(chan *frame.Notification)

	return c.startup()
}

// ssl 按配置中的严格程度开始TLS握手
func (c *Client) ssl() (err error) {
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

	switch c.dsn.SSL.Mode {
	case "prefer":
		if code == 'N' {
			// 服务器不支持，转用明文传输
			return nil
		} else if err = c.dsn.SSLCheck(); err != nil {
			// 服务器支持，但本地证书检查不通过，也转用明文传输
			return nil
		}
		tlsConfig.InsecureSkipVerify = true
	case "require":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		// 如果没有根证书，标记为只验证证书
		if _, err = os.Stat(c.dsn.SSL.RootCert); err == nil {
			tlsConfig.InsecureSkipVerify = true
		}
	case "verify-ca":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		tlsConfig.InsecureSkipVerify = false
	case "verify-full":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		tlsConfig.ServerName = c.dsn.Host
	default:
		return errors.New("pq: SSL unknown type")
	}

	if err = c.dsn.SSLCheck(); err != nil {
		return
	}
	tlsConfig.Renegotiation = tls.RenegotiateFreelyAsClient
	cert, err := tls.LoadX509KeyPair(c.dsn.SSL.Cert, c.dsn.SSL.Key)
	if err != nil {
		return err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	if _, err = os.Stat(c.dsn.SSL.RootCert); err == nil {
		tlsConfig.RootCAs = x509.NewCertPool()
		var raw []byte
		raw, err = ioutil.ReadFile(c.dsn.SSL.RootCert)
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

// startup 启动
func (c *Client) startup() (err error) {
	su := frame.NewStartup()
	for k, v := range c.dsn.Parameter {
		su.AddParam(k, v)
	}
	// 最后多补个空，表示结束
	su.WriteUint8(0)
	if err = c.writer.Fire(su.Data); err != nil {
		return
	}
	var out interface{}
	for {
		out, err = c.reader.Decode()
		if err != nil {
			return
		}
		switch f := out.(type) {
		case *frame.Error:
			f.Decode()
			log.Println(f.Error)
			return errors.New(f.Error.Error())
		case *frame.AuthRequest:
			switch f.GetType() {
			case frame.AuthTypePwd:
				ar := frame.NewAuthResponse()
				ar.Password(c.dsn.Password)
				if err = c.writer.Fire(ar.Data); err != nil {
					return
				}
			case frame.AuthTypeMd5:
				ar := frame.NewAuthResponse()
				ar.Md5Pwd(c.dsn.Parameter["user"], c.dsn.Password, string(f.GetMd5Salt()))
				if err = c.writer.Fire(ar.Data); err != nil {
					return
				}
			case frame.AuthTypeSASL:
				am := f.GetSASLAuthMechanism()
				if am != frame.AuthSASLSCRAMSHA256 && am != frame.AuthSASLSCRAMSHA256PLUS {
					return errors.New("不支持的SASL认证")
				}
				sc := scram2.NewClient(sha256.New, "", c.dsn.Password)
				sc.Step(nil)
				if sc.Err() != nil {
					return errors.New(fmt.Sprintf("SCRAM-SHA-256 error: %s", sc.Err().Error()))
				}

				//cli, err := scram.SHA256.NewClient("", c.dsn.Password, "")
				//if err != nil {
				//	return err
				//}
				//c.saslConversation = cli.NewConversation()
				//resp, err := c.saslConversation.Step("")
				//if err != nil {
				//	return err
				//}

				ar := frame.NewAuthSASLInitialResponse()
				ar.Mechanism(frame.AuthSASLSCRAMSHA256)
				ar.AuthResponse(string(sc.Out()))
				if err = c.writer.Fire(ar.Data); err != nil {
					return err
				}
			case frame.AuthTypeSASLContinue:
				resp, err := c.saslConversation.Step(string(f.GetSASLAuthData()))
				if err != nil {
					return err
				}
				ar := frame.NewAuthSASLResponse()
				ar.AuthResponse([]byte(resp))
				if err = c.writer.Fire(ar.Data); err != nil {
					return err
				}
			case frame.AuthTypeOk:
			}
		case *frame.ParameterStatus:
			f.Decode()
			c.parameterMaps[f.Name] = f.Value
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
			c.status = frame.TransactionStatus(f.Payload()[0])
			go c.readerLoop()
			return nil
		}
	}
}
