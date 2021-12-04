// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"context"
	"crypto/tls"
	"github.com/blusewang/pg/internal/frame"
	"github.com/blusewang/pg/internal/helper"
	"net"
	"time"
)

var (
	ListenMap = make(map[string]func(string))
)

type Client struct {
	ctx             context.Context         // 了解不深
	dsn             helper.DataSourceName   // 自动重连时需要
	conn            net.Conn                // 底层连接
	isTLS           *bool                   // 是否加密传输
	tlsConfig       *tls.Config             // tls 配置缓存
	writer          *frame.Encoder          // 流式编码
	reader          *frame.Decoder          // 流式解码
	backendPid      uint32                  // 业务过程中 后端PID 取消操作时需要
	backendKey      uint32                  // 业务过程中 后端口令 取消操作时需要
	parameterStatus map[string]string       // 服务器提供的属性参数
	Location        *time.Location          // 服务器端的时区
	status          frame.TransactionStatus // 业务状态
	frameChan       chan interface{}        // 帧流
	Err             *frame.Error            // 业务上最近的错误帧
	IOError         error                   // 底层通信错误
}

// Response 业务响应的数据集
type Response struct {
	ParameterDescription *frame.ParameterDescription
	Description          *frame.RowDescription
	DataRows             []*frame.DataRow
	Completion           *frame.CommandCompletion
}

// NewClient 创建客户端对象
func NewClient(ctx context.Context, dsn helper.DataSourceName) (c *Client, err error) {
	c = new(Client)
	c.dsn = dsn
	c.ctx = ctx
	c.status = frame.TransactionStatusNoReady
	c.parameterStatus = make(map[string]string)
	c.frameChan = make(chan interface{}, 100)
	return c, c.connect()
}

// Reconnect 重连
func (c *Client) Reconnect() (err error) {
	c.ctx = context.Background()
	return c.connect()
}
