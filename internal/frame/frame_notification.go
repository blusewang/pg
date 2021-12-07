// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type Notification struct {
	*Data
	Pid       uint32
	Condition string
	Text      string
}

func (n *Notification) Decode() {
	n.Pid = n.readUint32()
	n.Condition = n.readString()
	n.Text = n.readString()
}
