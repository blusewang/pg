// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

type PgResult struct {
}

// LastInsertId returns the database's auto-generated ID
// after, for example, an INSERT into a table with primary
// key.
func (pr *PgResult) LastInsertId() (int64, error) {
	return 0, nil
}

// RowsAffected returns the number of rows affected by the
// query.
func (pr *PgResult) RowsAffected() (int64, error) {
	return 0, nil
}
