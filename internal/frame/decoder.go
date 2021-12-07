// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(c net.Conn) *Decoder {
	return &Decoder{bufio.NewReader(c)}
}

func (d *Decoder) Decode() (out Frame, err error) {
	f := new(Data)
	f.Name, err = d.r.ReadByte()
	if err != nil {
		return
	}
	f.payload, err = d.r.Peek(4)
	if err != nil {
		return
	}
	f.length = binary.BigEndian.Uint32(f.payload)
	f.payload = make([]byte, f.length, f.length)
	_, err = io.ReadFull(d.r, f.payload)
	if err != nil {
		return
	}
	f.payload = f.payload[4:]
	out = detect(f)
	return
}

func detect(f *Data) (af Frame) {
	switch f.Name {
	case 'R':
		af = &AuthRequest{Data: f}
	case 'S':
		af = &ParameterStatus{Data: f}
	case 'K':
		af = &BackendKeyData{Data: f}
	case 'Z':
		af = &ReadyForQuery{Data: f}
	case '1':
		af = &ParseCompletion{Data: f}
	case 't':
		af = &ParameterDescription{Data: f}
	case 'T':
		af = &RowDescription{Data: f}
	case 'D':
		af = &DataRow{Data: f}
	case 'C':
		af = &CommandCompletion{Data: f}
	case 'E':
		af = &Error{Data: f}
	case 'A':
		af = &Notification{Data: f}
	case '3':
		af = &CloseComplete{Data: f}
	case 'I':
		af = &CloseComplete{Data: f}
	case 'n':
		af = &NoData{Data: f}
	case 'N':
		af = &NoticeResponse{Error{Data: f}}
	default:
		af = f
	}
	return
}
