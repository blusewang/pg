// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type errorStruct struct {
	LocalError       string `json:"local_error"`
	Fail             string `json:"error"`
	Code             string `json:"code"`
	Message          string `json:"message"`
	Detail           string `json:"detail"`
	Hint             string `json:"hint"`
	Position         string `json:"position"`
	InternalPosition string `json:"internal_position"`
	InternalQuery    string `json:"internal_query"`
	Where            string `json:"where"`
	Schema           string `json:"schema"`
	Table            string `json:"table"`
	Column           string `json:"Column"`
	DateType         string `json:"date_type"`
	Constraint       string `json:"constraint"`
	File             string `json:"file"`
	Line             string `json:"line"`
	Routine          string `json:"routine"`
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
