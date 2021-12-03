// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package client

import (
	"github.com/blusewang/pg/internal/frame"
	"io"
	"log"
)

func (c *Client) getFrames() (list []interface{}, err error) {
	for {
		out, isOpen := <-c.frameChan
		if !isOpen {
			err = io.EOF
			return
		}
		list = append(list, out)
		if e, has := out.(*frame.Error); has {
			err = e.Error
			return
		}
		if _, has := out.(*frame.ReadyForQuery); has {
			return
		}
	}
}

func (c *Client) QueryNoArgs(query string) (res Response, err error) {
	if err = c.writer.Encode(frame.NewSimpleQuery(query)); err != nil {
		return
	}
	fs, err := c.getFrames()
	if err != nil {
		return
	}
	for i := range fs {
		switch f := fs[i].(type) {
		case *frame.DataRow:
			f.Decode()
			res.DataRows = append(res.DataRows, f)
		case *frame.RowDescription:
			f.Decode()
			res.Description = f
		case *frame.CommandCompletion:
			res.Completion = f
		}
	}
	return
}

func (c *Client) Parse(name, query string) (res Response, err error) {
	c.Err = nil
	if err = c.writer.Encode(frame.NewParse(name, query)); err != nil {
		return
	}
	if err = c.writer.Encode(frame.NewDescribe(name)); err != nil {
		return
	}
	if err = c.writer.Encode(frame.NewSync()); err != nil {
		return
	}
	if err = c.writer.Flush(); err != nil {
		return
	}
	fs, err := c.getFrames()
	if err != nil {
		return
	}
	for i := range fs {
		switch f := fs[i].(type) {
		case *frame.RowDescription:
			f.Decode()
			res.Description = f
		}
	}
	if res.Completion == nil && c.Err != nil {
		err = c.Err.Error
	}
	return
}

// ParseExec BindFrame
func (c *Client) ParseExec(name string, args [][]byte) (res Response, err error) {
	c.Err = nil
	if err = c.writer.Encode(frame.NewBind(name, args)); err != nil {
		return
	}
	if err = c.writer.Encode(frame.NewExecute()); err != nil {
		return
	}
	if err = c.writer.Encode(frame.NewSync()); err != nil {
		return
	}
	if err = c.writer.Flush(); err != nil {
		return
	}
	fs, err := c.getFrames()
	if err != nil {
		return
	}
	for i := range fs {
		switch f := fs[i].(type) {
		case *frame.DataRow:
			f.Decode()
			res.DataRows = append(res.DataRows, f)
		case *frame.RowDescription:
			f.Decode()
			res.Description = f
		case *frame.CommandCompletion:
			res.Completion = f
		}
	}
	if res.Completion == nil && c.Err != nil {
		err = c.Err.Error
	}
	return
}

// ParseQuery BindFrame
func (c *Client) ParseQuery(name string, args [][]byte) (res Response, err error) {
	c.Err = nil
	if err = c.writer.Encode(frame.NewBind(name, args)); err != nil {
		return
	}
	if err = c.writer.Encode(frame.NewExecute()); err != nil {
		return
	}
	if err = c.writer.Encode(frame.NewSync()); err != nil {
		return
	}
	if err = c.writer.Flush(); err != nil {
		return
	}
	fs, err := c.getFrames()
	if err != nil {
		return
	}
	for i := range fs {
		switch f := fs[i].(type) {
		case *frame.DataRow:
			f.Decode()
			res.DataRows = append(res.DataRows, f)
		case *frame.RowDescription:
			f.Decode()
			res.Description = f
		case *frame.CommandCompletion:
			res.Completion = f
		}
	}
	if res.Completion == nil && c.Err != nil {
		err = c.Err.Error
	}
	return
}

func (c *Client) CloseParse(name string) (err error) {
	if err = c.writer.Encode(frame.NewCloseStat(name)); err != nil {
		return
	}
	if err = c.writer.Encode(frame.NewSync()); err != nil {
		return
	}
	if err = c.writer.Flush(); err != nil {
		return
	}
	_, err = c.getFrames()
	log.Println(err)
	return
}

func (c *Client) CancelRequest() (err error) {
	if err = c.writer.Encode(frame.NewCancelRequest(c.backendPid, c.backendKey)); err != nil {
		return
	}
	if err = c.writer.Encode(frame.NewSync()); err != nil {
		return
	}
	if err = c.writer.Flush(); err != nil {
		return
	}
	_, err = c.getFrames()
	return
}

func (c *Client) Terminate() (err error) {
	_ = c.writer.Fire(frame.NewTermination())
	close(c.responseChan)
	close(c.frameChan)
	c.parameterStatus = nil
	c.status = frame.TransactionStatusNoReady
	c.IOError = io.EOF
	return c.conn.Close()
}

func (c *Client) IsInTransaction() bool {
	return c.status == frame.TransactionStatusIdle || c.status == frame.TransactionStatusInFailedTransaction
}
