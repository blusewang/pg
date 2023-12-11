// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"context"
	"encoding/json"
	"github.com/blusewang/pg/internal/frame"
	"io"
	"log"
	"time"
)

// readerLoop 开启实时读取
// 实时处理异步数据
func (c *Client) readerLoop() {
	var out frame.Frame
	var err error
	for {
		out, err = c.reader.Decode()
		if err != nil {
			c.IOError = err
			log.Println(err)
			// 因为收到致命错误会立即退出此循环
			// 所以这里只会在网络断开或解码出错时触发

			if err == io.EOF {
				//close(c.frameChan)
				//close(c.NotifyChan)
				go func() {
					time.Sleep(time.Second * 2)
					log.Println("reconnect", c.connect(context.Background()))
				}()
			}
			return
		}
		switch f := out.(type) {
		case *frame.Notification:
			f.Decode()
			c.NotifyChan <- f
		case *frame.NoticeResponse:
			f.Decode()
			raw, _ := json.Marshal(f.Error.Error)
			log.Println(string(raw))
		case *frame.Error:
			f.Decode()
			c.Err = f
			c.status = frame.TransactionStatusNoReady
			raw, _ := json.Marshal(f.Error)
			log.Println(string(raw))
			c.frameChan <- f
			if f.Error.Fail == "FATAL" || f.Error.Fail == "PANIC" {
				// 这两种错误需立即断开连接
				//close(c.frameChan)
				//close(c.NotifyChan)
				go func() {
					time.Sleep(time.Second * 2)
					_ = c.Close()
					log.Println("reconnect", c.connect(context.Background()))
				}()
				return
			}
		case *frame.ReadyForQuery:
			c.status = frame.TransactionStatus(f.Payload()[0])
			c.frameChan <- f
		default:
			c.frameChan <- f
		}
	}
}

func (c *Client) getFrames() (list []frame.Frame, err error) {
	for {
		out, isOpen := <-c.frameChan
		if !isOpen {
			err = io.EOF
			return
		}
		list = append(list, out)
		if e, has := out.(*frame.Error); has {
			err = e.Error
		}
		if _, has := out.(*frame.ReadyForQuery); has {
			return
		}
	}
}
