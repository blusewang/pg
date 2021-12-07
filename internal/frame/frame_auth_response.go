// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

import (
	"crypto/md5"
	"fmt"
)

type AuthResponse struct {
	*Data
}

func NewAuthResponse() *AuthResponse {
	return &AuthResponse{&Data{
		Name:    'p',
		length:  0,
		payload: []byte{},
	}}
}

func (ar *AuthResponse) Password(str string) {
	ar.writeString(str)
}

func (ar *AuthResponse) Md5Pwd(user, pwd, salt string) {
	ar.writeString("md5" + ar.md5(ar.md5(pwd+user)+salt))
}

func (ar *AuthResponse) md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
