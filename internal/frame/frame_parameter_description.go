// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type ParameterDescription struct {
	*Data
	Count    uint16
	TypeOIDs []uint32
}

func (pd *ParameterDescription) Decode() {
	pd.Count = pd.readUint16()
	for i := uint16(0); i < pd.Count; i++ {
		pd.TypeOIDs = append(pd.TypeOIDs, pd.readUint32())
	}
	return
}
