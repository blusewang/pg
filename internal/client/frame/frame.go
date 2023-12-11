// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

import (
	"bytes"
	"encoding/binary"
)

type Data struct {
	Name     byte
	length   uint32
	position int
	payload  []byte
}

type Frame interface {
	Type() byte
	Length() uint32
	Payload() []byte
}

func (f *Data) Type() byte {
	return f.Name
}

func (f *Data) Length() uint32 {
	return f.length
}

func (f *Data) Payload() []byte {
	return f.payload
}

func (f *Data) writeUint32(uint uint32) {
	raw := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(raw, uint)
	f.payload = append(f.payload, raw...)
}

func (f *Data) writeUint16(uint uint16) {
	raw := []byte{0, 0}
	binary.BigEndian.PutUint16(raw, uint)
	f.payload = append(f.payload, raw...)
}

func (f *Data) WriteUint8(uint uint8) {
	f.payload = append(f.payload, uint)
}

func (f *Data) writeBytes(b []byte) {
	f.payload = append(f.payload, b...)
}

func (f *Data) writeString(str string) {
	str += string(byte(0))
	f.payload = append(f.payload, []byte(str)...)
}

func (f *Data) resetPosition() {
	f.position = 0
}

func (f *Data) readString() (str string) {
	if f.position >= len(f.payload) {
		return
	}
	position := bytes.IndexByte(f.payload[f.position:], 0)
	str = string(f.payload[f.position : f.position+position])
	f.position += position + 1
	return
}

func (f *Data) readUint64() (u uint64) {
	if f.position >= len(f.payload) {
		return
	}
	u = binary.BigEndian.Uint64(f.payload[f.position : f.position+8])
	f.position += 8
	return
}

func (f *Data) readUint32() (u uint32) {
	if f.position >= len(f.payload) {
		return
	}
	u = binary.BigEndian.Uint32(f.payload[f.position : f.position+4])
	f.position += 4
	return
}

func (f *Data) readUint16() (u uint16) {
	if f.position >= len(f.payload) {
		return
	}
	u = binary.BigEndian.Uint16(f.payload[f.position : f.position+2])
	f.position += 2
	return
}

func (f *Data) readUint6() (u uint8) {
	if f.position >= len(f.payload) {
		return
	}
	u = f.payload[f.position]
	f.position += 1
	return
}

func (f *Data) readLength(length int) (raw []byte) {
	if f.position+length > len(f.payload) {
		return []byte{}
	}
	raw = f.payload[f.position : f.position+length]
	f.position += length
	return raw
}

var _ Frame = new(Data)
