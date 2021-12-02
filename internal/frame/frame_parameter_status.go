// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type ParameterStatus struct {
	*Frame
	Name  string
	Value string
}

func (ps *ParameterStatus) Decode() {
	ps.resetPosition()
	ps.Name, ps.Value = ps.readString(), ps.readString()
	return
}
