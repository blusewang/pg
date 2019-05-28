// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
)

func (pi *PgIO) ssl() (err error) {
	code, err := pi.sslRequest()
	if err != nil {
		pi.IOError = err
		return
	}

	switch pi.dsn.SSL.Mode {
	case "prefer":
		if code == 'N' {
			return nil
		} else if err = pi.sslCheck(); err != nil {
			return nil
		}
		pi.tlsConfig.InsecureSkipVerify = true
	case "require":
		if code == 'N' {
			pi.IOError = errors.New("pq: SSL is not enabled on the server")
			return nil
		}
		// 如果没有根证书，标记为只验证证书
		if _, err = os.Stat(pi.dsn.SSL.RootCert); err == nil {
			pi.tlsConfig.InsecureSkipVerify = true
		}
	case "verify-ca":
		if code == 'N' {
			pi.IOError = errors.New("pq: SSL is not enabled on the server")
			return nil
		}
		pi.tlsConfig.InsecureSkipVerify = true
	case "verify-full":
		if code == 'N' {
			pi.IOError = errors.New("pq: SSL is not enabled on the server")
			return nil
		}
		pi.tlsConfig.ServerName = pi.dsn.Host

	default:
		return nil
	}

	if err = pi.sslCheck(); err != nil {
		return err
	}

	if err = pi.sslConfig(); err != nil {
		return err
	}

	pi.tlsConfig.Renegotiation = tls.RenegotiateFreelyAsClient

	pi.conn = tls.Client(pi.conn, &pi.tlsConfig)
	pi.reader = bufio.NewReader(pi.conn)

	return
}

func (pi *PgIO) sslRequest() (code byte, err error) {
	bs := NewPgMessage(IdentifiesSSLRequest)
	bs.addInt32(8)
	bs.addInt32(80877103)

	_ = bs.encode()
	if _, err = pi.conn.Write(bs.Content[4:]); err != nil {
		return
	}

	return pi.reader.ReadByte()
}

func (pi *PgIO) sslCheck() (err error) {
	if pi.dsn.SSL.Mode == "verify-ca" || pi.dsn.SSL.Mode == "verify-full" {
		if _, err = os.Stat(pi.dsn.SSL.RootCert); err != nil {
			return err
		}
	}

	if _, err = os.Stat(pi.dsn.SSL.Cert); err != nil {
		return err
	}

	info, err := os.Stat(pi.dsn.SSL.Key)
	if err != nil {
		return err
	}
	if info.Mode().Perm()&0077 != 0 {
		return errors.New("pq: Private key file has group or world access. Permissions should be u=rw (0600) or less")
	}

	return
}

func (pi *PgIO) sslConfig() (err error) {
	cert, err := tls.LoadX509KeyPair(pi.dsn.SSL.Cert, pi.dsn.SSL.Key)
	if err != nil {
		return err
	}
	pi.tlsConfig.Certificates = []tls.Certificate{cert}

	if _, err = os.Stat(pi.dsn.SSL.RootCert); err == nil {
		pi.tlsConfig.RootCAs = x509.NewCertPool()
		raw, err := ioutil.ReadFile(pi.dsn.SSL.RootCert)
		if err != nil {
			return err
		}
		if !pi.tlsConfig.RootCAs.AppendCertsFromPEM(raw) {
			return errors.New("pg: can't parse root cert")
		}
	}
	return nil
}
