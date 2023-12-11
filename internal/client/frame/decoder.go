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
	*bufio.Reader
}

func NewDecoder(c net.Conn) *Decoder {
	return &Decoder{bufio.NewReader(c)}
}

func (d *Decoder) Receive() (data *Data, err error) {
	data = new(Data)
	data.Name, err = d.ReadByte()
	if err != nil {
		return
	}
	data.payload, err = d.Peek(4)
	if err != nil {
		return
	}
	data.length = binary.BigEndian.Uint32(data.payload)
	data.payload = make([]byte, data.length)
	_, err = io.ReadFull(d, data.payload)
	if err != nil {
		return
	}
	data.payload = data.payload[4:]
	return data, err
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
