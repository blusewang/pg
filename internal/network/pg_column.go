// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

type PgColumn struct {
	Name         string
	TableOid     uint32
	Index        uint16
	TypeOid      uint32
	Len          uint16
	TypeModifier uint32
	Format       uint16
}
