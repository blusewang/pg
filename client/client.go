// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"context"
	"github.com/blusewang/pg/dsn"
	"github.com/blusewang/pg/internal/frame"
	"github.com/xdg-go/scram"
	"net"
	"time"
)

type Client struct {
	dsn              dsn.DataSourceName        // 取消任务时需要
	conn             net.Conn                  // 底层连接
	writer           *frame.Encoder            // 流式编码器
	reader           *frame.Decoder            // 流式解码器
	backendPid       uint32                    // 业务过程中 后端PID 取消操作时需要
	backendKey       uint32                    // 业务过程中 后端口令 取消操作时需要
	parameterMaps    map[string]string         // 服务器提供的属性参数
	StatementMaps    map[string]*Statement     // 句柄缓存，交给client管理，方便在断开时释放资源
	Location         *time.Location            // 服务器端的时区
	status           frame.TransactionStatus   // 业务状态
	frameChan        chan frame.Frame          // 帧流
	NotifyChan       chan *frame.Notification  // 异步消息通道
	Err              *frame.Error              // 业务上最近的错误帧
	IOError          error                     // 底层通信错误
	saslConversation *scram.ClientConversation // SASL
}

// Response 业务响应的数据集
type Response struct {
	ParameterDescription *frame.ParameterDescription
	Description          *frame.RowDescription
	DataRows             []*frame.DataRow
	Completion           *frame.CommandCompletion
}

// Statement 可缓存的语句
type Statement struct {
	Id       string
	SQL      string
	Response Response
}

// NewClient 创建客户端对象
func NewClient(ctx context.Context, dsn dsn.DataSourceName) (c *Client, err error) {
	c = new(Client)
	c.dsn = dsn
	c.status = frame.TransactionStatusNoReady
	c.parameterMaps = make(map[string]string)
	c.StatementMaps = make(map[string]*Statement)
	c.frameChan = make(chan frame.Frame)
	c.NotifyChan = make(chan *frame.Notification)
	return c, c.connect(ctx)
}
