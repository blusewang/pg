// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type CopyInResponse struct {
	*Data
	Format       CopyFormat
	ColumnLength uint16
}

func (cir *CopyInResponse) Decode() {
	cir.Format = CopyFormat(cir.readUint6())
	cir.ColumnLength = cir.readUint16()
}
