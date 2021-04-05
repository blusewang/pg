// Copyright 2020 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"log"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestNewPgConn(t *testing.T) {
	var ts = time.Duration(0)
	raw := []byte("-2 years 3 mons 12 days 19:00:00.2342")
	re := regexp.
		MustCompile(`^(?:([+-]?[0-9]+) years?)? ?(?:([+-]?[0-9]+) mons?)? ?(?:([+-]?[0-9]+) days?)? ?(?:([+-])?([0-9]+):([0-9]+):([0-9]+)(?:,|.([0-9]+))?)?$`)
	ps := re.FindStringSubmatch(string(raw))
	for i, p := range ps {
		log.Println(i, p)
	}
	// year
	i, _ := strconv.ParseInt(ps[1], 10, 64)
	ts += time.Duration(i*12*30*24) * time.Hour
	// month
	i, _ = strconv.ParseInt(ps[2], 10, 64)
	ts += time.Duration(i*30*24) * time.Hour
	// day
	i, _ = strconv.ParseInt(ps[3], 10, 64)
	ts += time.Duration(i*24) * time.Hour
	if ps[4] == "-" {
		// hour
		i, _ = strconv.ParseInt(ps[5], 10, 64)
		ts -= time.Duration(i) * time.Hour
		// min
		i, _ = strconv.ParseInt(ps[6], 10, 64)
		ts -= time.Duration(i) * time.Minute
		// sec
		i, _ = strconv.ParseInt(ps[7], 10, 64)
		ts -= time.Duration(i) * time.Second
		// ms
		i, _ = strconv.ParseInt(ps[8], 10, 64)
		ts -= time.Duration(i) * time.Millisecond
	} else {
		// hour
		i, _ = strconv.ParseInt(ps[5], 10, 64)
		ts += time.Duration(i) * time.Hour
		// min
		i, _ = strconv.ParseInt(ps[6], 10, 64)
		ts += time.Duration(i) * time.Minute
		// sec
		i, _ = strconv.ParseInt(ps[7], 10, 64)
		ts += time.Duration(i) * time.Second
		// ms
		i, _ = strconv.ParseInt(ps[8], 10, 64)
		ts += time.Duration(i) * time.Millisecond
	}
	log.Println(ts)
}
