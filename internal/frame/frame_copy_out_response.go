// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type CopyOutResponse struct {
	*Data
	Format        CopyFormat
	ColumnLength  uint16
	ColumnFormats []CopyFormat
}

type CopyFormat uint8

const (
	CopyFormatText   uint8 = 0
	CopyFormatBinary uint8 = 1
)

func (cor *CopyOutResponse) Decode() {
	cor.Format = CopyFormat(cor.readUint6())
	cor.ColumnLength = cor.readUint16()
	for i := uint16(0); i < cor.ColumnLength; i++ {
		cor.ColumnFormats = append(cor.ColumnFormats, CopyFormat(cor.readUint6()))
	}
}
