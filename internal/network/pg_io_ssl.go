// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import "log"

func (pi *PgIO) ssl() (err error) {
	switch pi.dsn.SSL.Mode {
	case "disable":
	case "allow":

	case "require":

	case "verify-ca":

	case "verify-full":

	case "prefer":
		fallthrough
	default:

	}
	return
}

func (pi *PgIO) sslRequest() (err error) {
	bs := NewPgMessage(IdentifiesSSLRequest)
	bs.addInt32(8)
	bs.addInt32(80877103)

	_ = bs.encode()
	_, err = pi.conn.Write(bs.Content)
	if err != nil {
		return
	}
	log.Println("send", bs.Content)

	id, err := pi.reader.ReadByte()
	if err != nil {
		pi.IOError = err
		return err
	}

	log.Println(string(id))

	if id != 'S' {

	}

	return
}
