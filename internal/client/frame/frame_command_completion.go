// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

import (
	"strconv"
	"strings"
)

type CommandCompletion struct {
	*Data
}

func (cc *CommandCompletion) Affected() (n int) {
	if cc == nil {
		return
	}
	arr := strings.Split(string(cc.payload[:len(cc.payload)-1]), " ")
	if len(arr) > 1 {
		n, _ = strconv.Atoi(arr[len(arr)-1])
	}
	return
}
