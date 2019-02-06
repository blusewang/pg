// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

type PgTx struct {
}

func (t *PgTx) Commit() error {
	return nil
}
func (t *PgTx) Rollback() error {
	return nil
}
