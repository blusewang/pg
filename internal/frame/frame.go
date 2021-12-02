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

type Frame struct {
	Name     byte
	length   uint32
	position int
	Payload  []byte
}

func (f *Frame) writeUint32(uint uint32) {
	raw := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(raw, uint)
	f.Payload = append(f.Payload, raw...)
}

func (f *Frame) writeUint16(uint uint16) {
	raw := []byte{0, 0}
	binary.BigEndian.PutUint16(raw, uint)
	f.Payload = append(f.Payload, raw...)
}

func (f *Frame) WriteUint8(uint uint8) {
	f.Payload = append(f.Payload, uint)
}

func (f *Frame) writeBytes(b []byte) {
	f.Payload = append(f.Payload, b...)
}

func (f *Frame) writeString(str string) {
	str += string(byte(0))
	f.Payload = append(f.Payload, []byte(str)...)
}

func (f *Frame) resetPosition() {
	f.position = 0
}

func (f *Frame) readString() (str string) {
	if f.position >= len(f.Payload) {
		return
	}
	position := bytes.IndexByte(f.Payload[f.position:], 0)
	str = string(f.Payload[f.position : f.position+position])
	f.position += position + 1
	return
}

func (f *Frame) readUint64() (u uint64) {
	if f.position >= len(f.Payload) {
		return
	}
	u = binary.BigEndian.Uint64(f.Payload[f.position : f.position+8])
	f.position += 8
	return
}

func (f *Frame) readUint32() (u uint32) {
	if f.position >= len(f.Payload) {
		return
	}
	u = binary.BigEndian.Uint32(f.Payload[f.position : f.position+4])
	f.position += 4
	return
}

func (f *Frame) readUint16() (u uint16) {
	if f.position >= len(f.Payload) {
		return
	}
	u = binary.BigEndian.Uint16(f.Payload[f.position : f.position+2])
	f.position += 2
	return
}

func (f *Frame) readUint6() (u uint8) {
	if f.position >= len(f.Payload) {
		return
	}
	u = f.Payload[f.position]
	f.position += 1
	return
}

func (f *Frame) readLength(length int) (raw []byte) {
	if f.position+length > len(f.Payload) {
		return []byte{}
	}
	raw = f.Payload[f.position : f.position+length]
	f.position += length
	return raw
}
