// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type Column struct {
	Name         string
	TableOid     uint32
	Index        uint16
	TypeOid      uint32
	Length       uint16
	TypeModifier uint32
	Format       uint16
}
type RowDescription struct {
	*Data
	Count   uint16
	Columns []Column
}

func (rd *RowDescription) Decode() {
	rd.resetPosition()
	rd.Count = rd.readUint16()
	rd.Columns = make([]Column, 0)
	for i := uint16(0); i < rd.Count; i++ {
		rd.Columns = append(rd.Columns, Column{
			Name:         rd.readString(),
			TableOid:     rd.readUint32(),
			Index:        rd.readUint16(),
			TypeOid:      rd.readUint32(),
			Length:       rd.readUint16(),
			TypeModifier: rd.readUint32(),
			Format:       rd.readUint16(),
		})
	}
}
