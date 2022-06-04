// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type AuthRequest struct {
	*Data
}

const (
	AuthTypeOk           uint32 = 0
	AuthTypePwd          uint32 = 3
	AuthTypeMd5          uint32 = 5
	AuthTypeSASL         uint32 = 10
	AuthTypeSASLContinue uint32 = 11
	AuthSASLSCRAMSHA256  string = "SCRAM-SHA-256"
)

func (ar *AuthRequest) GetType() uint32 {
	return ar.readUint32()
}

func (ar *AuthRequest) GetMd5Salt() []byte {
	return ar.readLength(4)
}

func (ar *AuthRequest) GetSASLAuthMechanism() string {
	return ar.readString()
}

func (ar *AuthRequest) GetSASLAuthData() []byte {
	return ar.payload[4:]
}
