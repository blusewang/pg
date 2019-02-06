// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package driver

import "database/sql/driver"

type PgDriverContext struct{}

func (dc *PgDriverContext) OpenConnector(name string) (driver.Connector, error) {
	return &PgConnector{}, nil
}
