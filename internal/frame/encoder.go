// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

import (
	"encoding/binary"
	"io"
	"net"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(c net.Conn) *Encoder {
	return &Encoder{c}
}

func (e *Encoder) Encode(f *Frame) (err error) {
	raw := make([]byte, 4)
	binary.BigEndian.PutUint32(raw, uint32(len(f.Payload))+4)
	if f.Name > 0 {
		raw = append([]byte{f.Name}, raw...)
	}
	_, err = e.w.Write(append(raw, f.Payload...))
	return
}
