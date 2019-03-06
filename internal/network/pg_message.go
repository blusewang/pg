// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func NewPgMessage(identifies Identifies) *PgMessage {
	msg := new(PgMessage)
	msg.Identifies = identifies
	msg.Len = 4
	msg.addInt32(0)
	return msg
}

type PgMessage struct {
	Identifies Identifies
	Len        uint32
	Content    []byte
	Position   uint32
}

func (pm *PgMessage) int32() (n uint32) {
	if pm.Position >= pm.Len {
		return 0
	}
	n = binary.BigEndian.Uint32(pm.Content[pm.Position:])
	defer pm.move(4)
	return
}

func (pm *PgMessage) int16() (n uint16) {
	if pm.Position >= pm.Len {
		return 0
	}
	n = binary.BigEndian.Uint16(pm.Content[pm.Position:])
	defer pm.move(2)
	return
}

func (pm *PgMessage) string() string {
	if pm.Position >= pm.Len {
		return ""
	}
	i := uint32(bytes.IndexByte(pm.Content[pm.Position:], 0))
	if i < 0 {
		return ""
	}
	defer pm.move(i + 1)
	return string(pm.Content[pm.Position : pm.Position+i])
}

func (pm *PgMessage) byte() byte {
	if pm.Position >= pm.Len {
		return 0
	}
	defer pm.move(1)
	return pm.Content[pm.Position]
}

func (pm *PgMessage) bytes(n uint32) []byte {

	if pm.Position >= pm.Len {
		return []byte{}
	}
	if pm.Position+n > pm.Len {
		n = pm.Len - pm.Position
	}
	defer pm.move(n)
	return pm.Content[pm.Position : pm.Position+n]
}

func (pm *PgMessage) move(n uint32) {
	pm.Position += n
}

func (pm *PgMessage) ParseError() (err *PgError) {
	if pm.Identifies != IdentifiesErrorResponse {
		return
	}
	err = new(PgError)
	for s := " "; s != ""; s = pm.string() {
		switch s[0] {
		case 'S':
			err.Severity = s[1:]
		case 'V':
			err.Text = s[1:]
		case 'C':
			err.Code, _ = strconv.Atoi(s[1:])
		case 'M':
			err.Message = s[1:]
		case 'D':
			err.Detail = s[1:]
		case 'H':
			err.Hint = s[1:]
		case 'P':
			err.Position, _ = strconv.Atoi(s[1:])
		case 'p':
			err.InternalPosition, _ = strconv.Atoi(s[1:])
		case 'q':
			err.InternalQuery = s[1:]
		case 'W':
			err.Where = s[1:]
		case 's':
			err.Schema = s[1:]
		case 't':
			err.Table = s[1:]
		case 'c':
			err.Column = s[1:]
		case 'd':
			err.DataType = s[1:]
		case 'n':
			err.Constraint = s[1:]
		case 'F':
			err.File = s[1:]
		case 'L':
			err.Line, _ = strconv.Atoi(s[1:])
		case 'R':
			err.Routine = s[1:]
		}
	}
	return
}

func (pm *PgMessage) columns() (list []PgColumn) {
	if pm.Identifies != IdentifiesRowDescription {
		return
	}
	count := pm.int16()
	for n := uint16(0); n < count; n++ {
		var c PgColumn
		c.Name = pm.string()
		c.TableOid = pm.int32()
		c.Index = pm.int16()
		c.TypeOid = pm.int32()
		c.Len = pm.int16()
		c.TypeModifier = pm.int32()
		c.Format = pm.int16()
		list = append(list, c)
	}
	return
}

func (pm *PgMessage) addInt32(n int) {
	x := make([]byte, 4)
	binary.BigEndian.PutUint32(x, uint32(n))
	pm.Content = append(pm.Content, x...)
}

func (pm *PgMessage) addInt16(n int) {
	x := make([]byte, 2)
	binary.BigEndian.PutUint16(x, uint16(n))
	pm.Content = append(pm.Content, x...)
}

func (pm *PgMessage) addString(s string) {
	pm.Content = append(pm.Content, s+"\000"...)
}

func (pm *PgMessage) addByte(c byte) {
	pm.Content = append(pm.Content, c)
}

func (pm *PgMessage) addBytes(v []byte) {
	pm.Content = append(pm.Content, v...)
}

func (pm *PgMessage) encode() (raw []byte) {
	raw = append(raw, byte(pm.Identifies))
	binary.BigEndian.PutUint32(pm.Content, uint32(len(pm.Content)))
	raw = append(raw, pm.Content...)
	return
}
