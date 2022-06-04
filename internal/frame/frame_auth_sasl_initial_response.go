// Copyright 2022 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

import "log"

type AuthSASLInitialResponse struct {
	*Data
}

func NewAuthSASLInitialResponse() *AuthSASLInitialResponse {
	return &AuthSASLInitialResponse{&Data{
		Name:    'p',
		length:  0,
		payload: []byte{},
	}}
}

func (ar *AuthSASLInitialResponse) Mechanism(str string) {
	ar.writeString(str)
}

func (ar *AuthSASLInitialResponse) AuthResponse(str string) {
	log.Println(uint8(len([]byte(str))), str)
	ar.writeUint32(uint32(len([]byte(str))))
	ar.writeBytes([]byte(str))
}
