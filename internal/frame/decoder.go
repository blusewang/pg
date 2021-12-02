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

func (d *Decoder) Decode() (out interface{}, err error) {
	f := new(Frame)
	f.Name, err = d.r.ReadByte()
	if err != nil {
		return
	}
	f.Payload, err = d.r.Peek(4)
	if err != nil {
		return
	}
	f.length = binary.BigEndian.Uint32(f.Payload)
	f.Payload = make([]byte, f.length, f.length)
	_, err = io.ReadFull(d.r, f.Payload)
	if err != nil {
		return
	}
	f.Payload = f.Payload[4:]
	out = detect(f)
	return
}

func detect(f *Frame) (af interface{}) {
	switch f.Name {
	case 'R':
		af = &AuthRequest{Frame: f}
	case 'S':
		af = &ParameterStatus{Frame: f}
	case 'K':
		af = &BackendKeyData{Frame: f}
	case 'Z':
		af = &ReadyForQuery{Frame: f}
	case '1':
		af = &ParseCompletion{Frame: f}
	case 't':
		af = &ParameterDescription{Frame: f}
	case 'T':
		af = &RowDescription{Frame: f}
	case 'D':
		af = &DataRow{Frame: f}
	case 'C':
		af = &CommandCompletion{Frame: f}
	case 'E':
		af = &Error{Frame: f}
	case 'A':
		af = &Notification{Frame: f}
	case '3':
		af = &CloseComplete{Frame: f}
	case 'I':
		af = &CloseComplete{Frame: f}
	case 'n':
		af = &NoData{Frame: f}
	case 'N':
		af = &NoticeResponse{Error{Frame: f}}
	default:
		af = f
	}
	return
}
