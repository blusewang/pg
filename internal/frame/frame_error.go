// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type errorStruct struct {
	LocalError       string `json:"local_error,omitempty"`
	Fail             string `json:"error,omitempty"`
	Code             string `json:"code,omitempty"`
	Message          string `json:"message,omitempty"`
	Detail           string `json:"detail,omitempty"`
	Hint             string `json:"hint,omitempty"`
	Position         string `json:"position,omitempty"`
	InternalPosition string `json:"internal_position,omitempty"`
	InternalQuery    string `json:"internal_query,omitempty"`
	Where            string `json:"where,omitempty"`
	Schema           string `json:"schema,omitempty"`
	Table            string `json:"table,omitempty"`
	Column           string `json:"Column,omitempty"`
	DateType         string `json:"date_type,omitempty"`
	Constraint       string `json:"constraint,omitempty"`
	File             string `json:"file,omitempty"`
	Line             string `json:"line,omitempty"`
	Routine          string `json:"routine,omitempty"`
}

func (e errorStruct) Error() string {
	return e.Message
}

type Error struct {
	*Frame
	Error errorStruct
}

func (e *Error) Decode() {
	for e.position < len(e.Payload)-1 {
		s := e.readString()
		switch s[0] {
		case 'S':
			e.Error.LocalError = s[1:]
		case 'V':
			e.Error.Fail = s[1:]
		case 'C':
			e.Error.Code = s[1:]
		case 'M':
			e.Error.Message = s[1:]
		case 'D':
			e.Error.Detail = s[1:]
		case 'H':
			e.Error.Hint = s[1:]
		case 'P':
			e.Error.Position = s[1:]
		case 'p':
			e.Error.InternalPosition = s[1:]
		case 'q':
			e.Error.InternalQuery = s[1:]
		case 'W':
			e.Error.Where = s[1:]
		case 's':
			e.Error.Schema = s[1:]
		case 't':
			e.Error.Table = s[1:]
		case 'c':
			e.Error.Column = s[1:]
		case 'd':
			e.Error.DateType = s[1:]
		case 'n':
			e.Error.Constraint = s[1:]
		case 'F':
			e.Error.File = s[1:]
		case 'L':
			e.Error.Line = s[1:]
		case 'R':
			e.Error.Routine = s[1:]
		}
	}
}
