// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type SSLRequest Data

func NewSSLRequest() *Data {
	f := &Data{
		Name:    0,
		payload: []byte{},
	}
	// hard code
	// uint32 = 80877103
	f.writeUint16(1234)
	f.writeUint16(5679)
	return f
}

func NewTermination() *Data {
	return &Data{
		Name:    'X',
		length:  4,
		payload: []byte{},
	}
}

func NewDescribe(stat string) *Data {
	var d = &Data{
		Name:    'D',
		payload: []byte{},
	}
	d.writeString("S" + stat)
	return d
}

func NewSync() *Data {
	return &Data{
		Name:    'S',
		length:  4,
		payload: []byte{},
	}
}

func NewParse(sid, query string) *Data {
	p := &Data{
		Name:    'P',
		payload: []byte{},
	}
	p.writeString(sid)
	p.writeString(query)
	p.writeUint16(0)
	return p
}

func NewExecute() *Data {
	e := &Data{
		Name:    'E',
		payload: []byte{},
	}
	// Portal
	e.writeString("")
	// returns
	e.writeUint32(0)
	return e
}

func NewSimpleQuery(stat string) *Data {
	e := &Data{
		Name:    'Q',
		payload: []byte{},
	}
	e.writeString(stat)
	return e
}

func NewCancelRequest(pid, key uint32) *Data {
	c := &Data{
		Name:    'F',
		payload: []byte{},
	}
	// hard code
	// uint32 = 80877102
	c.writeUint16(1234)
	c.writeUint16(5678)
	c.writeUint32(pid)
	c.writeUint32(key)
	return c
}

func NewCloseStat(stat string) *Data {
	c := &Data{
		Name:    'C',
		payload: []byte{},
	}
	c.WriteUint8('S')
	c.writeString(stat)
	return c
}

func NewFlush() *Data {
	return &Data{
		Name:    'H',
		payload: []byte{},
	}
}

type NoData struct {
	*Data
}

type EmptyQueryResponse struct {
	*Data
}

type CloseComplete struct {
	*Data
}

type ReadyForQuery struct {
	*Data
}

type ParseCompletion struct {
	*Data
}

type BindCompletion struct {
	*Data
}
