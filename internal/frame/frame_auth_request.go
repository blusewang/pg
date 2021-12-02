// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

import (
	"encoding/binary"
)

type AuthRequest struct {
	*Frame
}

const (
	AuthTypeOk  uint32 = 0
	AuthTypePwd uint32 = 3
	AuthTypeMd5 uint32 = 5
)

func (ar *AuthRequest) GetType() uint32 {
	return binary.BigEndian.Uint32(ar.Payload[0:4])
}

func (ar *AuthRequest) GetMd5Salt() []byte {
	return ar.Payload[4:]
}
