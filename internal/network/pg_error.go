// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import (
	"fmt"
)

type PgError struct {
	Severity         string `json:"severity"`
	Text             string `json:"text"`
	Code             int    `json:"code"`
	Message          string `json:"message"`
	Detail           string `json:"detail"`
	Hint             string `json:"hint"`
	Position         int    `json:"position"`
	InternalPosition int    `json:"internal_position"`
	InternalQuery    string `json:"internal_query"`
	Where            string `json:"where"`
	Schema           string `json:"schema"`
	Table            string `json:"table"`
	Column           string `json:"column"`
	DataType         string `json:"data_type"`
	Constraint       string `json:"constraint"`
	File             string `json:"file"`
	Line             int    `json:"line"`
	Routine          string `json:"routine"`
}

func (e *PgError) Error() string {
	return fmt.Sprintf("pg %v %v", e.Severity, e.Message)
}
