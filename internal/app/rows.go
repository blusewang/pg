package app

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/blusewang/pg/v2/internal/client/frame"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Rows struct {
	location *time.Location
	columns  *frame.RowDescription
	rows     []*frame.DataRow
	position int
}

func (r *Rows) HasNextResultSet() bool {
	return false
}

func (r *Rows) NextResultSet() error {
	return nil
}

func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	switch frame.PgType(r.columns.Columns[index].TypeOid) {
	case frame.PgTypeBool:
		return reflect.TypeOf((*bool)(nil)).Elem()
	case frame.PgTypeDate, frame.PgTypeTime, frame.PgTypeTimestamp, frame.PgTypeTimestamptz, frame.PgTypeTimetz:
		return reflect.TypeOf((*time.Time)(nil)).Elem()
	case frame.PgTypeInt2, frame.PgTypeInt4, frame.PgTypeInt8:
		return reflect.TypeOf((*int64)(nil)).Elem()
	case frame.PgTypeFloat4, frame.PgTypeFloat8, frame.PgTypeNumeric:
		return reflect.TypeOf((*float64)(nil)).Elem()
	case frame.PgTypeText, frame.PgTypeVarchar, frame.PgTypeChar, frame.PgTypeUuid:
		return reflect.TypeOf((*string)(nil)).Elem()
	case frame.PgTypePoint:
		//return reflect.TypeOf("")
		return reflect.TypeOf((*[]float64)(nil)).Elem()
	case frame.PgTypeJson, frame.PgTypeJsonb:
		return reflect.TypeOf((*json.RawMessage)(nil)).Elem()
	case frame.PgTypeBytea:
		return reflect.TypeOf((*[]byte)(nil)).Elem()

	case frame.PgTypeArrBool:
		return reflect.TypeOf((*[]bool)(nil)).Elem()
	case frame.PgTypeArrDate, frame.PgTypeArrTime, frame.PgTypeArrTimestamp, frame.PgTypeArrTimestamptz, frame.PgTypeArrTimetz:
		return reflect.TypeOf((*[]time.Time)(nil)).Elem()
	case frame.PgTypeArrInt2, frame.PgTypeArrInt4, frame.PgTypeArrInt8:
		return reflect.TypeOf((*[]int64)(nil)).Elem()
	case frame.PgTypeArrFloat4, frame.PgTypeArrFloat8, frame.PgTypeArrNumeric:
		return reflect.TypeOf((*[]float64)(nil)).Elem()
	case frame.PgTypeArrVarchar, frame.PgTypeArrChar, frame.PgTypeArrText, frame.PgTypeArrUuid, frame.PgTypeArrJson, frame.PgTypeArrJsonb:
		return reflect.TypeOf((*[]string)(nil)).Elem()
	case frame.PgTypeArrBytea:
		return reflect.TypeOf((*[][]byte)(nil)).Elem()
	case frame.PgTypeRecord:
		return reflect.TypeOf((*[]string)(nil)).Elem()

	default:
		return reflect.TypeOf((*string)(nil)).Elem()
	}
}

func (r *Rows) ColumnTypePrecisionScale(_ int) (precision, scale int64, ok bool) {
	return 0, 0, false
}

func (r *Rows) ColumnTypeNullable(_ int) (nullable, ok bool) {
	return true, false
}

func (r *Rows) ColumnTypeLength(index int) (length int64, ok bool) {
	return int64(r.columns.Columns[index].Length), true
}

func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	return frame.PgTypeMap[frame.PgType(r.columns.Columns[index].TypeOid)]
}

func (r *Rows) Columns() (cols []string) {
	if r.columns == nil {
		return
	}
	for _, v := range r.columns.Columns {
		cols = append(cols, v.Name)
	}
	return
}

func (r *Rows) Close() error {
	r.position = 0
	r.rows = nil
	r.columns = nil
	return nil
}

func (r *Rows) Next(dest []driver.Value) error {
	var rowsLen = len(r.rows)
	if rowsLen == 0 {
		return sql.ErrNoRows
	} else if r.position < 0 || r.position >= rowsLen {
		return fmt.Errorf("pg_rows rows length is %v but position is %v", rowsLen, r.position)
	} else if r.position == rowsLen {
		return fmt.Errorf("pg_rows rows length is %v but position is %v", rowsLen, r.position)
	}
	for i, v := range r.rows[r.position].DataArr {
		dest[i] = r.data2Value(v, r.columns.Columns[i])
	}
	r.position += 1
	return nil
}

func (r *Rows) data2Value(raw []byte, col frame.Column) interface{} {
	if raw == nil {
		return nil
	}
	switch col.TypeOid {
	case frame.PgTypeBool:
		return string(raw)[0] == 't'
	case frame.PgTypeText, frame.PgTypeChar, frame.PgTypeVarchar:
		return string(raw)

	case frame.PgTypeTimestamptz:
		t, _ := time.Parse("2006-01-02 15:04:05-07", string(raw))
		return t.In(r.location)
	case frame.PgTypeTimestamp:
		t, _ := time.Parse("2006-01-02 15:04:05", string(raw))
		return t.In(r.location).Add(-8 * time.Hour)
	case frame.PgTypeDate:
		t, _ := time.Parse(time.DateOnly, string(raw))
		return t.In(r.location).Add(-8 * time.Hour)
	case frame.PgTypeTime:
		t, _ := time.Parse(time.TimeOnly, string(raw))
		return t.In(r.location).Add(-8 * time.Hour)
	case frame.PgTypeTimetz:
		t, _ := time.Parse("15:04:05-07", string(raw))
		return t.In(r.location)

	case frame.PgTypeInt2, frame.PgTypeInt4, frame.PgTypeInt8:
		var i, _ = strconv.Atoi(string(raw))
		return i
	case frame.PgTypeFloat4, frame.PgTypeFloat8, frame.PgTypeNumeric:
		var f, _ = strconv.ParseFloat(string(raw), 64)
		return f
	case frame.PgTypeUuid:
		return string(raw)
	case frame.PgTypePoint:
		raw = bytes.ReplaceAll(raw, []byte("("), []byte("["))
		raw = bytes.ReplaceAll(raw, []byte(")"), []byte("]"))
		var arr []float64
		_ = json.Unmarshal(raw, &arr)
		return arr
		//return string(raw)
	case frame.PgTypeJson, frame.PgTypeJsonb:
		return json.RawMessage(raw)
	case frame.PgTypeArrBool:
		raw = bytes.ReplaceAll(raw, []byte("{"), []byte("["))
		raw = bytes.ReplaceAll(raw, []byte("}"), []byte("]"))
		raw = bytes.ReplaceAll(raw, []byte("t"), []byte("true"))
		raw = bytes.ReplaceAll(raw, []byte("f"), []byte("false"))
		var arr []int64
		_ = json.Unmarshal(raw, &arr)
		return arr
	case frame.PgTypeArrInt4, frame.PgTypeArrInt8:
		var str = string(raw)
		var arr []int64
		if strings.HasPrefix(str, "{") && len(str) > 2 {
			str = str[1 : len(str)-1]
			for _, v := range strings.Split(str, ",") {
				if v != "" {
					n, _ := strconv.ParseInt(v, 0, 64)
					arr = append(arr, n)
				}
			}
		}
		return arr
	case frame.PgTypeArrFloat4, frame.PgTypeArrFloat8, frame.PgTypeArrNumeric:
		raw = bytes.ReplaceAll(raw, []byte("{"), []byte("["))
		raw = bytes.ReplaceAll(raw, []byte("}"), []byte("]"))
		if bytes.HasPrefix(raw, []byte("[[[")) {
			var arr [][][]float64
			_ = json.Unmarshal(raw, &arr)
			return arr
		} else if bytes.HasPrefix(raw, []byte("[[")) {
			var arr [][]float64
			_ = json.Unmarshal(raw, &arr)
			return arr
		} else if bytes.HasPrefix(raw, []byte("[")) {
			var arr []float64
			_ = json.Unmarshal(raw, &arr)
			return arr
		}
		return []float64{}
	case frame.PgTypeArrText, frame.PgTypeArrChar, frame.PgTypeArrVarchar, frame.PgTypeArrUuid:
		var ss = pgStringArr{Raw: bytes.Runes(raw)}
		ss.parse()
		return ss.rs
	case frame.PgTypeRecord:
		return strings.Split(string(raw[1:len(raw)-1]), ",")
	default:
		return string(raw)
	}
}

type pgStringArr struct {
	Raw      []rune
	position int
	buf      []rune
	rs       []string
}

func (sa *pgStringArr) parse() {
	var length = len(sa.Raw)
	if length < 3 {
		return
	}
	if sa.Raw[0] != '{' {
		return
	}
	sa.position = 1
	for {
		if sa.position >= length-1 {
			break
		}
		if sa.Raw[sa.position] == '"' {
			sa.position += 1
			sa.readUtil('"')
			var str = string(sa.buf)
			str = strings.Replace(str, `\'`, `'`, -1)
			str = strings.Replace(str, `\"`, `"`, -1)
			sa.rs = append(sa.rs, str)
		} else {
			sa.readUtil(',')
			sa.rs = append(sa.rs, string(sa.buf))
		}
	}
}

func (sa *pgStringArr) readUtil(r rune) {
	sa.buf = []rune{}
	for {
		var rr = sa.Raw[sa.position]
		sa.position++
		if rr == r && sa.Raw[sa.position-2] != '\\' {
			break
		} else if rr == '}' && sa.Raw[sa.position-2] != '\\' {
			break
		} else {
			sa.buf = append(sa.buf, rr)
		}
	}
}

var _ driver.Rows = new(Rows)
var _ driver.RowsColumnTypeDatabaseTypeName = new(Rows)
var _ driver.RowsColumnTypeLength = new(Rows)
var _ driver.RowsColumnTypeNullable = new(Rows)
var _ driver.RowsColumnTypePrecisionScale = new(Rows)
var _ driver.RowsColumnTypeScanType = new(Rows)
var _ driver.RowsNextResultSet = new(Rows)
