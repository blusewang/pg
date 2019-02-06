// Copyright 2019 MQ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package network

import "log"

type TransactionStatus byte

const (
	TransactionStatusIdle                = 'I'
	TransactionStatusIdleInTransaction   = 'T'
	TransactionStatusInFailedTransaction = 'E'
)

func (t TransactionStatus) String() string {
	switch t {
	case TransactionStatusIdle:
		return "idle"
	case TransactionStatusIdleInTransaction:
		return "idle in transaction"
	case TransactionStatusInFailedTransaction:
		return "in a failed transaction"
	default:
		log.Panicf("unknown transactionStatus %d", t)
		return ""
	}
}
