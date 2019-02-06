// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import (
	"fmt"
	"strconv"
)

type PgError struct {
	Severity string `json:"severity"`
	Text     string `json:"text"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Routine  string `json:"routine"`
}

func (e *PgError) Error() string {
	return fmt.Sprintf("pg %v %v", e.Severity, e.Message)
}

func ParsePgError(buf *ReadBuffer) (err *PgError) {
	err = &PgError{}
	s := " "
	for ; s != ""; s = buf.string() {
		switch s[0] {
		case 'S':
			err.Severity = s[1:]
		case 'V':
			err.Text = s[1:]
		case 'C':
			err.Code, _ = strconv.Atoi(s[1:])
		case 'M':
			err.Message = s[1:]
		case 'F':
			err.File = s[1:]
		case 'L':
			err.Line, _ = strconv.Atoi(s[1:])
		case 'R':
			err.Routine = s[1:]
		}
	}
	return
}
