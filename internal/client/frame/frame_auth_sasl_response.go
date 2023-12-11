// Copyright 2022 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type AuthSASLResponse struct {
	*Data
}

func NewAuthSASLResponse() *AuthSASLResponse {
	return &AuthSASLResponse{&Data{
		Name:    'p',
		length:  0,
		payload: []byte{},
	}}
}

func (ar *AuthSASLResponse) AuthResponse(raw []byte) {
	ar.writeBytes(raw)
}
