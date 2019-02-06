// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import (
	"log"
	"testing"
)

func TestWriteBufferString(t *testing.T) {
	b := WriteBuffer{}
	b.int32(0)
	b.string("abcd")
	log.Println(b.buf, []byte{255}, []uint8{255, 255})
	//var bin []byte
	//bin = append(bin,[]byte("sdf")...,byte(0))
}
