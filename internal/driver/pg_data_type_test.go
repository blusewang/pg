// Copyright 2022 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestArray(t *testing.T) {
	var arr [][]bool
	str := "{{f,f,t,f}}"
	str = strings.Replace(str, "{", "[", -1)
	str = strings.Replace(str, "}", "]", -1)
	str = strings.Replace(str, "f", "false", -1)
	str = strings.Replace(str, "t", "true", -1)
	log.Println(str)
	log.Println(json.Unmarshal([]byte(str), &arr))

	log.Println(arr)
}
