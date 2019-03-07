// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/blusewang/pg/internal/network"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 针对需改造的数据做转换
func convert(raw []byte, col network.PgColumn, fieldLen uint32, location *time.Location) driver.Value {
	if fieldLen == 4294967295 {
		// is nil
		//return nil
		//log.Println(raw,col.Name,col.TypeOid,col.Len)
		//return []byte{}
	}
	if col.Format != 0 {
		panic("not support binary data")
	}

	switch PgType(col.TypeOid) {
	case PgTypeBool:
		return string(raw)[0] == 't'
	case PgTypeText, PgTypeChar, PgTypeVarchar:
		return string(raw)
	case PgTypeBytea:
		var b, _ = parseBytea(raw)
		return b
	case PgTypeTimestamptz:
		return parseTs(location, string(raw))
	case PgTypeTimestamp, PgTypeDate:
		return parseTs(nil, string(raw))
	case PgTypeTime:
		return mustParse("15:04:05", PgTypeTime, raw)
	case PgTypeTimetz:
		return mustParse("15:04:05-07", PgTypeTimetz, raw)
	case PgTypeInt2, PgTypeInt4, PgTypeInt8:
		var i, _ = strconv.Atoi(string(raw))
		return i
	case PgTypeFloat4, PgTypeFloat8, PgTypeNumeric:
		var f, _ = strconv.ParseFloat(string(raw), 64)
		return f
	case PgTypeJson, PgTypeJsonb:
		return string(raw)
	case PgTypeArrInt4:
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
	case PgTypeArrFloat4, PgTypeArrFloat8:
		var str = string(raw)
		var arr []float64
		if strings.HasPrefix(str, "{") && len(str) > 2 {
			str = str[1 : len(str)-1]
			for _, v := range strings.Split(str, ",") {
				if v != "" {
					n, _ := strconv.ParseFloat(v, 64)
					arr = append(arr, n)
				}
			}
		}
		return arr
	case PgTypeArrText, PgTypeArrChar, PgTypeArrVarchar:
		var str = string(raw)
		var arr []string
		if strings.HasPrefix(str, "{") && len(str) > 2 {
			str = str[2 : len(str)-2]
			for _, v := range strings.Split(str, "','") {
				if v != "" {
					arr = append(arr, v)
				}
			}
		}
		return arr
	default:
		return raw
	}
}

// Parse a bytea value received from the server.  Both "hex" and the legacy
// "escape" format are supported.
func parseBytea(s []byte) (result []byte, err error) {
	if len(s) >= 2 && bytes.Equal(s[:2], []byte("\\x")) {
		// bytea_output = hex
		s = s[2:] // trim off leading "\\x"
		result = make([]byte, hex.DecodedLen(len(s)))
		_, err := hex.Decode(result, s)
		if err != nil {
			return nil, err
		}
	} else {
		// bytea_output = escape
		for len(s) > 0 {
			if s[0] == '\\' {
				// escaped '\\'
				if len(s) >= 2 && s[1] == '\\' {
					result = append(result, '\\')
					s = s[2:]
					continue
				}

				// '\\' followed by an octal number
				if len(s) < 4 {
					return nil, fmt.Errorf("invalid bytea sequence %v", s)
				}
				r, err := strconv.ParseInt(string(s[1:4]), 8, 9)
				if err != nil {
					return nil, fmt.Errorf("could not parse bytea value: %s", err.Error())
				}
				result = append(result, byte(r))
				s = s[4:]
			} else {
				// We hit an unescaped, raw byte.  Try to read in as many as
				// possible in one go.
				i := bytes.IndexByte(s, '\\')
				if i == -1 {
					result = append(result, s...)
					break
				}
				result = append(result, s[:i]...)
				s = s[i:]
			}
		}
	}

	return result, nil
}

var infinityTsEnabled = false
var infinityTsNegative time.Time
var infinityTsPositive time.Time
var errInvalidTimestamp = errors.New("invalid timestamp")

func mustParse(f string, typ PgType, s []byte) time.Time {
	str := string(s)

	// check for a 30-minute-offset timezone
	if (typ == PgTypeTimestamptz || typ == PgTypeTimetz) &&
		str[len(str)-3] == ':' {
		f += ":00"
	}
	t, _ := time.Parse(f, str)
	return t
}

// This is a time function specific to the Postgres default DateStyle
// setting ("ISO, MDY"), the only one we currently support. This
// accounts for the discrepancies between the parsing available with
// time.Parse and the Postgres date formatting quirks.
func parseTs(currentLocation *time.Location, str string) interface{} {
	switch str {
	case "-infinity":
		if infinityTsEnabled {
			return infinityTsNegative
		}
		return []byte(str)
	case "infinity":
		if infinityTsEnabled {
			return infinityTsPositive
		}
		return []byte(str)
	}
	t, err := ParseTimestamp(currentLocation, str)
	if err != nil {
		return nil
	}
	return t
}

// ParseTimestamp parses Postgres' text format. It returns a time.Time in
// currentLocation iff that time's offset agrees with the offset sent from the
// Postgres server. Otherwise, ParseTimestamp returns a time.Time with the
// fixed offset offset provided by the Postgres server.
func ParseTimestamp(currentLocation *time.Location, str string) (time.Time, error) {
	p := timestampParser{}

	monSep := strings.IndexRune(str, '-')
	// this is Gregorian year, not ISO Year
	// In Gregorian system, the year 1 BC is followed by AD 1
	year := p.mustAtoi(str, 0, monSep)
	daySep := monSep + 3
	month := p.mustAtoi(str, monSep+1, daySep)
	p.expect(str, '-', daySep)
	timeSep := daySep + 3
	day := p.mustAtoi(str, daySep+1, timeSep)

	minLen := monSep + len("01-01") + 1

	isBC := strings.HasSuffix(str, " BC")
	if isBC {
		minLen += 3
	}

	var hour, minute, second int
	if len(str) > minLen {
		p.expect(str, ' ', timeSep)
		minSep := timeSep + 3
		p.expect(str, ':', minSep)
		hour = p.mustAtoi(str, timeSep+1, minSep)
		secSep := minSep + 3
		p.expect(str, ':', secSep)
		minute = p.mustAtoi(str, minSep+1, secSep)
		secEnd := secSep + 3
		second = p.mustAtoi(str, secSep+1, secEnd)
	}
	remainderIdx := monSep + len("01-01 00:00:00") + 1
	// Three optional (but ordered) sections follow: the
	// fractional seconds, the time zone offset, and the BC
	// designation. We set them up here and adjust the other
	// offsets if the preceding sections exist.

	nanoSec := 0
	tzOff := 0

	if remainderIdx < len(str) && str[remainderIdx] == '.' {
		fracStart := remainderIdx + 1
		fracOff := strings.IndexAny(str[fracStart:], "-+ ")
		if fracOff < 0 {
			fracOff = len(str) - fracStart
		}
		fracSec := p.mustAtoi(str, fracStart, fracStart+fracOff)
		nanoSec = fracSec * (1000000000 / int(math.Pow(10, float64(fracOff))))

		remainderIdx += fracOff + 1
	}
	if tzStart := remainderIdx; tzStart < len(str) && (str[tzStart] == '-' || str[tzStart] == '+') {
		// time zone separator is always '-' or '+' (UTC is +00)
		var tzSign int
		switch c := str[tzStart]; c {
		case '-':
			tzSign = -1
		case '+':
			tzSign = +1
		default:
			return time.Time{}, fmt.Errorf("expected '-' or '+' at position %v; got %v", tzStart, c)
		}
		tzHours := p.mustAtoi(str, tzStart+1, tzStart+3)
		remainderIdx += 3
		var tzMin, tzSec int
		if remainderIdx < len(str) && str[remainderIdx] == ':' {
			tzMin = p.mustAtoi(str, remainderIdx+1, remainderIdx+3)
			remainderIdx += 3
		}
		if remainderIdx < len(str) && str[remainderIdx] == ':' {
			tzSec = p.mustAtoi(str, remainderIdx+1, remainderIdx+3)
			remainderIdx += 3
		}
		tzOff = tzSign * ((tzHours * 60 * 60) + (tzMin * 60) + tzSec)
	}
	var isoYear int

	if isBC {
		isoYear = 1 - year
		remainderIdx += 3
	} else {
		isoYear = year
	}
	if remainderIdx < len(str) {
		return time.Time{}, fmt.Errorf("expected end of input, got %v", str[remainderIdx:])
	}
	t := time.Date(isoYear, time.Month(month), day,
		hour, minute, second, nanoSec,
		globalLocationCache.getLocation(tzOff))

	if currentLocation != nil {
		// Set the location of the returned Time based on the session's
		// TimeZone value, but only if the local time zone database agrees with
		// the remote database on the offset.
		lt := t.In(currentLocation)
		_, newOff := lt.Zone()
		if newOff == tzOff {
			t = lt
		}
	}

	return t, p.err
}

type timestampParser struct {
	err error
}

func (p *timestampParser) expect(str string, char byte, pos int) {
	if p.err != nil {
		return
	}
	if pos+1 > len(str) {
		p.err = errInvalidTimestamp
		return
	}
	if c := str[pos]; c != char && p.err == nil {
		p.err = fmt.Errorf("expected '%v' at position %v; got '%v'", char, pos, c)
	}
}

func (p *timestampParser) mustAtoi(str string, begin int, end int) int {
	if p.err != nil {
		return 0
	}
	if begin < 0 || end < 0 || begin > end || end > len(str) {
		p.err = errInvalidTimestamp
		return 0
	}
	result, err := strconv.Atoi(str[begin:end])
	if err != nil {
		if p.err == nil {
			p.err = fmt.Errorf("expected number; got '%v'", str)
		}
		return 0
	}
	return result
}

// The location cache caches the time zones typically used by the client.
type locationCache struct {
	cache map[int]*time.Location
	lock  sync.Mutex
}

// All connections share the same list of timezones. Benchmarking shows that
// about 5% speed could be gained by putting the cache in the connection and
// losing the mutex, at the cost of a small amount of memory and a somewhat
// significant increase in code complexity.
var globalLocationCache = newLocationCache()

func newLocationCache() *locationCache {
	return &locationCache{cache: make(map[int]*time.Location)}
}

// Returns the cached timezone for the specified offset, creating and caching
// it if necessary.
func (c *locationCache) getLocation(offset int) *time.Location {
	c.lock.Lock()
	defer c.lock.Unlock()

	location, ok := c.cache[offset]
	if !ok {
		location = time.FixedZone("", offset)
		c.cache[offset] = location
	}

	return location
}
