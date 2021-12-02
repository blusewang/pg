// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"context"
	"github.com/blusewang/pg/internal/frame"
	"github.com/blusewang/pg/internal/helper"
	"net"
	"time"
)

var (
	ListenMap = make(map[string]func(string))
)

type Client struct {
	ctx             context.Context
	conn            net.Conn
	writer          *frame.Encoder
	reader          *frame.Decoder
	backendPid      uint32
	backendKey      uint32
	parameterStatus map[string]string
	Location        *time.Location
	status          frame.TransactionStatus
	frameChan       chan interface{}
	responseChan    chan Response
	Err             *frame.Error
	IOError         error
}

type Response struct {
	Description *frame.RowDescription
	DataRows    []*frame.DataRow
	Completion  *frame.CommandCompletion
}

func NewClient(ctx context.Context, dsn helper.DataSourceName) (c *Client, err error) {
	c = new(Client)
	c.status = frame.TransactionStatusNoReady
	c.parameterStatus = make(map[string]string)
	c.frameChan = make(chan interface{})
	c.responseChan = make(chan Response)
	return c, c.connect(ctx, dsn)
}

func (c *Client) IsAlive() bool {
	return c.IOError != nil
}
