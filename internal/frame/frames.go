// Copyright 2021 YBCZ, Inc. All rights reserved.
//
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file in the root of the source
// tree.

package frame

type SSLRequest Frame

func NewSSLRequest() *Frame {
	f := &Frame{
		Name:    0,
		Payload: []byte{},
	}
	// hard code
	// uint32 = 80877103
	f.writeUint16(1234)
	f.writeUint16(5679)
	return f
}

func NewTermination() *Frame {
	return &Frame{
		Name:    'X',
		length:  4,
		Payload: []byte{},
	}
}

func NewDescribe(stat string) *Frame {
	var d = &Frame{
		Name:    'D',
		Payload: []byte{},
	}
	d.writeString("S" + stat)
	return d
}

func NewSync() *Frame {
	return &Frame{
		Name:    'S',
		length:  4,
		Payload: []byte{},
	}
}

func NewParse(sid, query string) *Frame {
	p := &Frame{
		Name:    'P',
		Payload: []byte{},
	}
	p.writeString(sid)
	p.writeString(query)
	p.writeUint16(0)
	return p
}

func NewBind(stat string, args [][]byte) *Frame {
	b := &Frame{
		Name:    'B',
		Payload: []byte{},
	}
	// Portal
	b.writeString("")
	// statement
	b.writeString(stat)
	// parameter formats
	b.writeUint16(0)
	// parameter values
	b.writeUint16(uint16(len(args)))
	for _, arg := range args {
		// length
		b.writeUint32(uint32(len(arg)))
		// data
		b.writeBytes(arg)
	}
	// Result formats
	b.writeUint16(0)
	return b
}

func NewExecute() *Frame {
	e := &Frame{
		Name:    'E',
		Payload: []byte{},
	}
	// Portal
	e.writeString("")
	// returns
	e.writeUint32(0)
	return e
}

func NewSimpleQuery(stat string) *Frame {
	e := &Frame{
		Name:    'Q',
		Payload: []byte{},
	}
	e.writeString(stat)
	return e
}

func NewCancelRequest(pid, key uint32) *Frame {
	c := &Frame{
		Name:    'F',
		Payload: []byte{},
	}
	// hard code
	// uint32 = 80877102
	c.writeUint16(1234)
	c.writeUint16(5678)
	c.writeUint32(pid)
	c.writeUint32(key)
	return c
}

func NewCloseStat(stat string) *Frame {
	c := &Frame{
		Name:    'C',
		Payload: []byte{},
	}
	c.WriteUint8('S')
	c.writeString(stat)
	return c
}

func NewFlush() *Frame {
	return &Frame{
		Name:    'H',
		Payload: []byte{},
	}
}

type NoData struct {
	*Frame
}

type EmptyQueryResponse struct {
	*Frame
}

type CloseComplete struct {
	*Frame
}

type ReadyForQuery struct {
	*Frame
}

type ParseCompletion struct {
	*Frame
}

type BindCompletion struct {
	*Frame
}
