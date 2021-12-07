// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type Startup struct {
	*Data
}

func NewStartup() *Startup {
	sf := &Startup{&Data{
		Name:    0,
		length:  0,
		payload: []byte{}},
	}
	// v3.0
	sf.writeUint16(3)
	sf.writeUint16(0)
	return sf
}

func (sf *Startup) AddParam(name, value string) {
	sf.writeString(name)
	sf.writeString(value)
}
