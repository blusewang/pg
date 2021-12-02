// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type BackendKeyData struct {
	*Frame
	Pid uint32
	Key uint32
}

func (bkd *BackendKeyData) Decode() {
	bkd.Pid = bkd.readUint32()
	bkd.Key = bkd.readUint32()
}
