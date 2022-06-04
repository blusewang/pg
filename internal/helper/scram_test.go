// Copyright 2022 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package helper

import (
	"github.com/xdg-go/scram"
	"log"
	"testing"
)

func TestScram(t *testing.T) {
	c, err := scram.SHA256.NewClient("", "abc123", "")
	if err != nil {
		t.Fatal(err)
	}
	conv := c.NewConversation()
	log.Println(conv.Step(""))
	log.Println(conv.Step("r=JmHACioiuoEbR0thl4AYbZDza66674kMEcCCElaoWoygORgl,s=HFtp2H6LfYfbEbgPN6S4ww==,i=4096"))
}
