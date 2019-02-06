// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import (
	"bytes"
	"encoding/binary"
)

type Oid uint32

type ReadBuffer []byte

func (b *ReadBuffer) int32() (n int) {
	n = int(int32(binary.BigEndian.Uint32(*b)))
	*b = (*b)[4:]
	return
}

func (b *ReadBuffer) oid() (n Oid) {
	n = Oid(binary.BigEndian.Uint32(*b))
	*b = (*b)[4:]
	return
}

// N.B: this is actually an unsigned 16-bit integer, unlike int32
func (b *ReadBuffer) int16() (n int) {
	n = int(binary.BigEndian.Uint16(*b))
	*b = (*b)[2:]
	return
}

func (b *ReadBuffer) string() string {
	i := bytes.IndexByte(*b, 0)
	if i < 0 {
		return ""
	}
	s := (*b)[:i]
	*b = (*b)[i+1:]
	return string(s)
}

func (b *ReadBuffer) next(n int) (v []byte) {
	v = (*b)[:n]
	*b = (*b)[n:]
	return
}

func (b *ReadBuffer) byte() byte {
	return b.next(1)[0]
}

type WriteBuffer struct {
	buf []byte
	pos int
}

func (b *WriteBuffer) int32(n int) {
	x := make([]byte, 4)
	binary.BigEndian.PutUint32(x, uint32(n))
	b.buf = append(b.buf, x...)
}

func (b *WriteBuffer) int16(n int) {
	x := make([]byte, 2)
	binary.BigEndian.PutUint16(x, uint16(n))
	b.buf = append(b.buf, x...)
}

func (b *WriteBuffer) string(s string) {
	b.buf = append(b.buf, s+"\000"...)
}

func (b *WriteBuffer) byte(c byte) {
	b.buf = append(b.buf, c)
}

func (b *WriteBuffer) bytes(v []byte) {
	b.buf = append(b.buf, v...)
}

func (b *WriteBuffer) wrap() []byte {
	p := b.buf[b.pos:]
	binary.BigEndian.PutUint32(p, uint32(len(p)))
	return b.buf
}

func (b *WriteBuffer) next(c byte) {
	p := b.buf[b.pos:]
	binary.BigEndian.PutUint32(p, uint32(len(p)))
	b.pos = len(b.buf) + 1
	b.buf = append(b.buf, c, 0, 0, 0, 0)
}
