package client

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/blusewang/pg/v2/internal/client/frame"
	"github.com/blusewang/pg/v2/internal/client/scram"
	"net"
	"os"
	"time"
)

func NewClient() (c *Client) {
	c = new(Client)
	c.ConnectStatus = ConnectStatusConnecting
	c.status = frame.TransactionStatusNoReady
	c.parameterMaps = make(map[string]string)
	return
}

type NotificationHandler func(pid uint32, channel, message string)

type ParseResponse struct {
	Parameters *frame.ParameterDescription
	Rows       *frame.RowDescription
}

type BindExecResponse struct {
	DataRows   []*frame.DataRow
	Completion *frame.CommandCompletion
}

type SimpleQueryResponse struct {
	ParseResponse
	BindExecResponse
}

type ConnectStatus int

const (
	ConnectStatusConnecting = iota
	ConnectStatusConnected
	ConnectStatusDisconnected
)

type Client struct {
	cn                  net.Conn                // TCP 连接
	ctx                 context.Context         // 初始上下文
	Dsn                 DataSourceName          // 数据源
	writer              *frame.Encoder          // 流式编码器
	reader              *frame.Decoder          // 流式解码器
	backendPid          uint32                  // 业务过程中 后端PID 取消操作时需要
	backendKey          uint32                  // 业务过程中 后端口令 取消操作时需要
	parameterMaps       map[string]string       // 服务器提供的属性参数
	Location            *time.Location          // 服务器端的时区
	status              frame.TransactionStatus // 业务状态
	ConnectStatus       ConnectStatus           // 连接状态
	notificationHandler NotificationHandler     // Listen 消息
}

func (c *Client) Connect(ctx context.Context, dsn DataSourceName) (err error) {
	c.ctx = ctx
	c.Dsn = dsn
	nw, addr, timeout := dsn.Address()
	c.ConnectStatus = ConnectStatusConnecting
	d := net.Dialer{Timeout: timeout}
	c.cn, err = d.DialContext(ctx, nw, addr)
	if err != nil {
		return
	}
	c.writer = frame.NewEncoder(c.cn)
	c.reader = frame.NewDecoder(c.cn)
	return
}

// AutoSSL 按配置中的严格程度开始TLS握手
func (c *Client) AutoSSL() (err error) {
	if c.Dsn.SSL.Mode == "disable" || c.Dsn.SSL.Mode == "allow" {
		return
	}
	var tlsConfig tls.Config
	if err = c.writer.Send(frame.NewSSLRequest()); err != nil {
		return
	}
	var code byte
	code, err = c.reader.ReadByte()
	if err != nil {
		return
	}

	switch c.Dsn.SSL.Mode {
	case "prefer":
		if code == 'N' {
			// 服务器不支持，转用明文传输
			return nil
		} else if err = c.Dsn.SSLCheck(); err != nil {
			// 服务器支持，但本地证书检查不通过，也转用明文传输
			return nil
		}
		tlsConfig.InsecureSkipVerify = true
	case "require":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		// 如果没有根证书，标记为只验证证书
		if _, err = os.Stat(c.Dsn.SSL.RootCert); err == nil {
			tlsConfig.InsecureSkipVerify = true
		}
	case "verify-ca":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		tlsConfig.InsecureSkipVerify = false
	case "verify-full":
		if code == 'N' {
			return errors.New("pq: SSL is not enabled on the server")
		}
		tlsConfig.ServerName = c.Dsn.Host
	default:
		return errors.New("pq: SSL unknown type")
	}

	if err = c.Dsn.SSLCheck(); err != nil {
		return
	}
	tlsConfig.Renegotiation = tls.RenegotiateFreelyAsClient
	cert, err := tls.LoadX509KeyPair(c.Dsn.SSL.Cert, c.Dsn.SSL.Key)
	if err != nil {
		return err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	if _, err = os.Stat(c.Dsn.SSL.RootCert); err == nil {
		tlsConfig.RootCAs = x509.NewCertPool()
		var raw []byte
		raw, err = os.ReadFile(c.Dsn.SSL.RootCert)
		if err != nil {
			return err
		}
		if !tlsConfig.RootCAs.AppendCertsFromPEM(raw) {
			return errors.New("pg: can't parse root cert")
		}
	}

	// 升级至TLS
	c.cn = tls.Client(c.cn, &tlsConfig)
	c.writer = frame.NewEncoder(c.cn)
	c.reader = frame.NewDecoder(c.cn)
	return
}

// Startup 启动
func (c *Client) Startup() (err error) {
	// 发启动指令
	su := frame.NewStartup()
	for k, v := range c.Dsn.Parameter {
		su.AddParam(k, v)
	}
	// 最后多补个空，表示结束
	su.WriteUint8(0)
	if err = c.writer.Send(su.Data); err != nil {
		return
	}

	// 认证
	sc := scram.NewClient(sha256.New, c.Dsn.Parameter["user"], c.Dsn.Password)
	for {
		d, ioErr := c.reader.Receive()
		if ioErr != nil {
			return c.handleIOError(ioErr)
		}
		switch d.Type() {
		case frame.TypeError:
			return c.handlePgError(d)
		case frame.TypeAuthRequest:
			auth := frame.AuthRequest{Data: d}
			switch auth.GetType() {
			case frame.AuthTypePwd:
				ar := frame.NewAuthResponse()
				ar.Password(c.Dsn.Password)
				if err = c.writer.Send(ar.Data); err != nil {
					return
				}
			case frame.AuthTypeMd5:
				ar := frame.NewAuthResponse()
				ar.Md5Pwd("", c.Dsn.Password, string(auth.GetMd5Salt()))
				if err = c.writer.Send(ar.Data); err != nil {
					return
				}
			case frame.AuthTypeSASL:
				am := auth.GetSASLAuthMechanism()
				if am != frame.AuthSASLSCRAMSHA256 && am != frame.AuthSASLSCRAMSHA256PLUS {
					return errors.New("不支持的SASL认证")
				}
				sc.Step(nil)
				if sc.Err() != nil {
					return errors.New(fmt.Sprintf("SCRAM-SHA-256 error: %s", sc.Err().Error()))
				}

				ar := frame.NewAuthSASLInitialResponse()
				ar.Mechanism(frame.AuthSASLSCRAMSHA256)
				ar.AuthResponse(string(sc.Out()))
				if err = c.writer.Send(ar.Data); err != nil {
					return err
				}
			case frame.AuthTypeSASLContinue:
				sc.Step(auth.GetSASLAuthData())
				if sc.Err() != nil {
					return errors.New(fmt.Sprintf("SCRAM-SHA-256 error: %s", sc.Err().Error()))
				}
				ar := frame.NewAuthSASLResponse()
				ar.AuthResponse(sc.Out())
				if err = c.writer.Send(ar.Data); err != nil {
					return err
				}
			case frame.AuthTypeOk:
			}
		case frame.TypeParameterStatus:
			p := frame.ParameterStatus{Data: d}
			p.Decode()
			c.parameterMaps[p.Name] = p.Value
			if p.Name == "TimeZone" {
				c.Location, err = time.LoadLocation(p.Value)
				if err != nil {
					c.Location = nil
				}
			}
		case frame.TypeBackendKeyData:
			bkd := frame.BackendKeyData{Data: d}
			bkd.Decode()
			c.backendPid = bkd.Pid
			c.backendKey = bkd.Key
		case frame.TypeReadyForQuery:
			c.status = frame.TransactionStatus(d.Payload()[0])
			c.ConnectStatus = ConnectStatusConnected
			return nil
		}
	}
}

func (c *Client) QueryNoArgs(query string) (res SimpleQueryResponse, err error) {
	if err = c.writer.Send(frame.NewSimpleQuery(query)); err != nil {
		return res, c.handleIOError(err)
	}

	for {
		f, ioErr := c.reader.Receive()
		if ioErr != nil {
			return res, c.handleIOError(ioErr)
		}
		switch f.Type() {
		case frame.TypeDataRow:
			d := frame.DataRow{Data: f}
			d.Decode()
			res.DataRows = append(res.DataRows, &d)
		case frame.TypeRowDescription:
			d := frame.RowDescription{Data: f}
			d.Decode()
			res.Rows = &d
		case frame.TypeCommandCompletion:
			d := frame.CommandCompletion{Data: f}
			res.Completion = &d
		case frame.TypeReadyForQuery:
			c.status = frame.TransactionStatus(f.Payload()[0])
			return
		case frame.TypeError:
			err = c.handlePgError(f)
		}
	}
}

func (c *Client) Parse(name, query string) (res ParseResponse, err error) {
	if err = c.writer.Buff(frame.NewParse(name, query)); err != nil {
		return res, c.handleIOError(err)
	}
	if err = c.writer.Buff(frame.NewDescribe(name)); err != nil {
		return res, c.handleIOError(err)
	}
	if err = c.writer.Buff(frame.NewSync()); err != nil {
		return res, c.handleIOError(err)
	}
	if err = c.writer.Flush(); err != nil {
		return res, c.handleIOError(err)
	}

	for {
		f, ioErr := c.reader.Receive()
		if ioErr != nil {
			return res, c.handleIOError(ioErr)
		}

		switch f.Type() {
		case frame.TypeParameterDescription:
			d := frame.ParameterDescription{Data: f}
			d.Decode()
			res.Parameters = &d

		case frame.TypeRowDescription:
			d := frame.RowDescription{Data: f}
			d.Decode()
			res.Rows = &d
		case frame.TypeReadyForQuery:
			c.status = frame.TransactionStatus(f.Payload()[0])
			return
		case frame.TypeError:
			err = c.handlePgError(f)
		}
	}
}

func (c *Client) BindExec(name string, args []driver.Value) (res BindExecResponse, err error) {
	if err = c.writer.Buff(frame.NewBind(name, args)); err != nil {
		return res, c.handleIOError(err)
	}
	if err = c.writer.Buff(frame.NewExecute()); err != nil {
		return res, c.handleIOError(err)
	}
	if err = c.writer.Buff(frame.NewSync()); err != nil {
		return res, c.handleIOError(err)
	}
	if err = c.writer.Flush(); err != nil {
		return res, c.handleIOError(err)
	}

	for {
		f, ioErr := c.reader.Receive()
		if ioErr != nil {
			return res, c.handleIOError(ioErr)
		}
		switch f.Type() {
		case frame.TypeNoticeResponse:
			d := frame.NoticeResponse{Error: frame.Error{Data: f}}
			d.Decode()
			d.Log()
		case frame.TypeDataRow:
			d := frame.DataRow{Data: f}
			d.Decode()
			res.DataRows = append(res.DataRows, &d)
		case frame.TypeCommandCompletion:
			d := frame.CommandCompletion{Data: f}
			res.Completion = &d
		case frame.TypeReadyForQuery:
			c.status = frame.TransactionStatus(f.Payload()[0])
			return
		case frame.TypeError:
			err = c.handlePgError(f)
		}
	}
}

func (c *Client) CloseParse(name string) (err error) {
	if err = c.writer.Buff(frame.NewCloseStat(name)); err != nil {
		return c.handleIOError(err)
	}
	if err = c.writer.Buff(frame.NewSync()); err != nil {
		return c.handleIOError(err)
	}
	if err = c.writer.Flush(); err != nil {
		return c.handleIOError(err)
	}
	for {
		d, ioErr := c.reader.Receive()
		if ioErr != nil {
			return c.handleIOError(ioErr)
		}
		switch d.Type() {
		case frame.TypeReadyForQuery:
			c.status = frame.TransactionStatus(d.Payload()[0])
			return
		case frame.TypeError:
			err = c.handlePgError(d)
		}
	}
}

func (c *Client) Listen(channel string) (err error) {
	_, err = c.QueryNoArgs("listen " + channel)
	return
}

func (c *Client) GetNotification() (pid uint32, channel, message string, err error) {
	d, ioErr := c.reader.Receive()
	if ioErr != nil {
		return
	}
	for {
		switch d.Type() {
		case frame.TypeNotification:
			n := frame.Notification{Data: d}
			n.Decode()
			return n.Pid, n.Condition, n.Text, nil
		default:

		}
	}
}

// CancelRequest 建立新连接，使用PID+口令从新连接中发出指令
func (c *Client) CancelRequest() (err error) {
	// TODO 建立新连接，使用PID+口令从新连接中发出指令
	return
}

func (c *Client) Terminate() (err error) {
	_ = c.writer.Send(frame.NewTermination())
	defer func() {
		c.ConnectStatus = ConnectStatusDisconnected
	}()
	return c.cn.Close()
}

func (c *Client) CloseConn() (err error) {
	return c.cn.Close()
}

func (c *Client) IsInTransaction() bool {
	return c.status == frame.TransactionStatusIdleInTransaction || c.status == frame.TransactionStatusInFailedTransaction
}

func (c *Client) handleIOError(err error) error {
	c.ConnectStatus = ConnectStatusDisconnected
	_ = c.cn.Close()
	return err
}

func (c *Client) handlePgError(d *frame.Data) error {
	e := frame.Error{Data: d}
	e.Decode()
	c.status = frame.TransactionStatusNoReady
	if e.Error.Fail == "FATAL" || e.Error.Fail == "PANIC" {
		// 这两种错误需立即断开连接
		go func() {
			_ = c.cn.Close()
			c.ConnectStatus = ConnectStatusDisconnected
		}()
	}
	return e.Error
}
