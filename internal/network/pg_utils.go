// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import (
	"go/types"
	"strconv"
	"time"
)

//   int64
//   float64
//   bool
//   []byte
//   string
//   time.Time
func value2bytes(value interface{}) []byte {
	switch value.(type) {
	case types.Nil:
		return nil
	case int64:
		return strconv.AppendInt(nil, value.(int64), 10)
	case float64:
		return strconv.AppendFloat(nil, value.(float64), 'f', -1, 64)
	case bool:
		return strconv.AppendBool(nil, value.(bool))
	case []byte:
		return value.([]byte)
	case string:
		return []byte(value.(string))
	case time.Time:
		return formatTimestamp(value.(time.Time))
	case *time.Time:
		return formatTimestamp(*value.(*time.Time))
	default:
		return []byte{}
	}
}

func formatTimestamp(t time.Time) []byte {
	// Need to send dates before 0001 A.D. with " BC" suffix, instead of the
	// minus sign preferred by Go.
	// Beware, "0000" in ISO is "1 BC", "-0001" is "2 BC" and so on
	bc := false
	if t.Year() <= 0 {
		// flip year sign, and add 1, e.g: "0" will be "1", and "-10" will be "11"
		t = t.AddDate((-t.Year())*2+1, 0, 0)
		bc = true
	}
	b := []byte(t.Format("2006-01-02 15:04:05.999999999Z07:00"))

	_, offset := t.Zone()
	offset = offset % 60
	if offset != 0 {
		// RFC3339Nano already printed the minus sign
		if offset < 0 {
			offset = -offset
		}

		b = append(b, ':')
		if offset < 10 {
			b = append(b, '0')
		}
		b = strconv.AppendInt(b, int64(offset), 10)
	}

	if bc {
		b = append(b, " BC"...)
	}
	return b
}
