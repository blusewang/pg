// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type DataRow struct {
	*Frame
	Count   uint16
	DataArr [][]byte
}

func (dr *DataRow) Decode() {
	dr.resetPosition()
	dr.Count = dr.readUint16()
	for i := uint16(0); i < dr.Count; i++ {
		length := dr.readUint32()
		if length == MaxUint32 {
			dr.DataArr = append(dr.DataArr, nil)
		} else {
			dr.DataArr = append(dr.DataArr, dr.readLength(int(length)))
		}
	}
}
