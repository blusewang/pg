// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

import (
	"database/sql/driver"
	"go/types"
	"math"
	"strconv"
	"time"
)

func NewBind(stat string, args []driver.Value) *Data {
	b := &Data{
		Name:    'B',
		payload: []byte{},
	}
	// Portal
	b.writeString("")
	// statement
	b.writeString(stat)
	// parameter formats
	b.writeUint16(0)
	// parameter values
	b.writeUint16(uint16(len(args)))
	for _, arg := range args {
		if arg == nil {
			b.writeUint32(math.MaxUint32)
		} else {
			// 转换
			raw := value2Row(arg)
			// length
			b.writeUint32(uint32(len(raw)))
			// data
			b.writeBytes(raw)
		}
	}
	// Result formats
	b.writeUint16(0)
	return b
}

func value2Row(value driver.Value) []byte {
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
		return []byte(value.(time.Time).Format(time.RFC3339Nano))
	case *time.Time:
		return []byte(value.(*time.Time).Format(time.RFC3339Nano))
	default:
		return []byte{}
	}
}
