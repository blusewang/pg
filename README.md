# pg

用于golang database/sql 的PostgreSQL驱动

[![GoDoc](https://godoc.org/github.com/blusewang/pg?status.svg)](https://godoc.org/github.com/blusewang/pg)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/blusewang/pg/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/blusewang/pg)](https://goreportcard.com/report/github.com/blusewang/pg)

#### Go Version Support

![Go version](https://img.shields.io/badge/Go-1.20x-brightgreen.svg)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fblusewang%2Fpg.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fblusewang%2Fpg?ref=badge_shield)

#### PostgreSQL Version Support

![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-10.5-brightgreen.svg)
![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-11.4-brightgreen.svg)
![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-12.0-brightgreen.svg)
![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-13.0-brightgreen.svg)
![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-14.0-brightgreen.svg)
![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-15.0-brightgreen.svg)
![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-16.0-brightgreen.svg)

## 特性

* Scan() 时允许传入指针，完美对应数据库中的`null`
* 所有查询全部自动`prepare`并缓存。
* 支持`pg://`前缀的URI
* 支持`Listen`式的异步消息订阅

## 安装

	go get github.com/blusewang/pg/v2

## 使用

### 通用查询

```golang

sql.Register("pg", pg.Driver{})
db, err := sql.Open("pg", "pg://user:password@dbhost.yourdomain.com/database_name?application_name=app_name&sslmode=verify-full")
if err != nil {
    return err
}
defer db.Close()
rows, err := db.Query("select * from table_name where id>$1", 0)
if err != nil {
    return err
}
...

```

### 异步订阅

```golang
    l, err := pg.NewListener("pg://user:password@dbhost.yourdomain.com/database_name?application_name=app_name&sslmode=verify-full")
    if err != nil {
        return
    }
    
    if err := l.Listen("channel_name")); err != nil {
        return        
    }
    go func() {
        time.Sleep(time.Minute)
        _ = l.Terminate()
    }()
    for {
        pid, channel, message, err := l.GetNotification()
        if err == io.EOF {
            return err
        } else if err != nil {
            return err
        }
        log.Println(pid, channel, message)
    }
```

## 文档

更多的细节及使用示例，参见： <https://pkg.go.dev/github.com/blusewang/pg/v2>.


### 配置要求

* 配置上需将`sql.SetMaxIdleConns(x)`、`sql.SetMaxOpenConns(x)`两处的x设置为相同的值，才能让缓存实现价值。

## 协议实现

- 本驱动是为Web服务而设计

| 状态                       | 功能        | 备注                            |
|:-------------------------|:----------|:------------------------------|
| <ul><li>- [x] </li></ul> | 启动        | 支持：无密码、明文密码、md5、SCRAM-SHA-256 |
| <ul><li>- [x] </li></ul> | 简单查询      | 必备                            |
| <ul><li>- [x] </li></ul> | 扩展查询      | 必备                            |
| <ul><li>- [x] </li></ul> | 取消正在处理的请求 | 必备                            |
| <ul><li>- [x] </li></ul> | 终止        | 必备                            |
| <ul><li>- [x] </li></ul> | SSL会话加密   | 远程安全                          |
| <ul><li>- [x] </li></ul> | 异步        | listen/notify                 |

暂不支持`SCRAM-SHA-256-PLUS`。当`pg_hba.conf`使用：
```conf
hostssl all             all             0.0.0.0/0               scram-sha-256 clientcert=verify-full
```
时，Server会发类型为`SCRAM-SHA-256-PLUS`类型的请求。此时客户端将自动降级至`SCRAM-SHA-256`。
## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fblusewang%2Fpg.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fblusewang%2Fpg?ref=badge_large)