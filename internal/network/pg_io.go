// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func NewPgIO() *PgIO {
	pi := &PgIO{}
	pi.serverConf = make(map[string]string)
	return pi
}

type PgIO struct {
	net.Conn
	reader     *bufio.Reader
	txStatus   TransactionStatus
	serverPid  int
	serverConf map[string]string
	backendKey int
	location   *time.Location
}

func (pi *PgIO) Dial(network, address string, timeout time.Duration) (err error) {
	pi.Conn, err = net.DialTimeout(network, address, timeout)
	if err == nil {
		pi.reader = bufio.NewReader(pi.Conn)
	}
	return
}

func (pi *PgIO) DialContext(context context.Context, network, address string, timeout time.Duration) (err error) {
	d := net.Dialer{Timeout: timeout}
	pi.Conn, err = d.DialContext(context, network, address)
	if err == nil {
		pi.reader = bufio.NewReader(pi.Conn)
	}
	return
}

func (pi *PgIO) StartUp(p map[string]string, pwd string) (err error) {
	var wb WriteBuffer
	wb.int32(0)
	wb.int32(196608)
	for k, v := range p {
		wb.string(k)
		wb.string(v)
	}
	wb.byte(0)
	_, err = pi.Write(wb.wrap())
	if err != nil {
		return
	}

	for {
		t, reader, err := pi.receive()
		if err != nil {
			return err
		}
		switch t {
		case 'R':
			err = pi.auth(reader, p["user"], pwd)
			if err != nil {
				return err
			}
		case 'S':
			k := reader.string()
			v := reader.string()
			if k == "TimeZone" {
				pi.location, err = time.LoadLocation(v)
				if err != nil {
					err = nil
					pi.location = nil
				}
			}
			log.Println(string(t), k, v)
			pi.serverConf[k] = v
		case 'K':
			pi.serverPid = reader.int32()
			pi.backendKey = reader.int32()
		case 'Z':
			pi.txStatus = TransactionStatus(reader.byte())
			return nil
		case 'E':
			err = ParsePgError(reader)
			log.Println(err)
			return err
		default:
			return fmt.Errorf("unknown response for startup: %q,%v", t, reader.string())
		}
	}
}

func (pi *PgIO) auth(reader *ReadBuffer, user, password string) (err error) {

	switch code := reader.int32(); code {
	case 0:
		// OK
		break
	case 3:
		// 明文密码
		msg := pi.newMsgByCommand('p')
		msg.string(password)
		_, err = pi.Write(msg.wrap())
		if err != nil {
			return err
		}
		t, buf, err := pi.receive()
		if err != nil {
			return err
		}

		if t == 'E' {
			return ParsePgError(buf)
		}

		if t == 'R' && buf.int32() != 0 {
			return fmt.Errorf("unexpected authentication response: %q", t)
		}
	case 5:
		// MD5密码
		msg := pi.newMsgByCommand('p')
		msg.string("md5" + pi.md5(pi.md5(password+user)+string(reader.next(4))))

		_, err = pi.Write(msg.wrap())
		if err != nil {
			return err
		}
		t, buf, err := pi.receive()
		if err != nil {
			return err
		}

		if t == 'E' {
			return ParsePgError(buf)
		}

		if t == 'R' && buf.int32() != 0 {
			return fmt.Errorf("unexpected authentication response: %q", t)
		}
	}
	return
}

func (pi *PgIO) md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (pi *PgIO) newMsgByCommand(b byte) *WriteBuffer {
	r := &WriteBuffer{pos: 1}
	r.byte(b)
	r.int32(0)
	return r
}

func (pi *PgIO) receive() (t byte, buf *ReadBuffer, err error) {
	for {
		buf = &ReadBuffer{}
		t, err = pi.reader.ReadByte()
		if err != nil {
			return t, buf, err
		}

		raw, err := pi.reader.Peek(4)
		if err != nil {
			return t, buf, err
		}
		n := int(binary.BigEndian.Uint32(raw))

		raw = make([]byte, n)
		_, err = io.ReadFull(pi.reader, raw)
		if err != nil {
			return t, buf, err
		}

		*buf = raw[4:]

		if t == 'N' {
			// ignore
			log.Println("pg-io received a notice message ->", string(t), buf.string())
		} else {
			return t, buf, err
		}
	}
}
