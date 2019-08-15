# pg
用于golang database/sql 的PostgreSQL驱动

[![GoDoc](https://godoc.org/github.com/blusewang/pg?status.svg)](https://godoc.org/github.com/blusewang/pg)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/blusewang/pg/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/blusewang/pg)](https://goreportcard.com/report/github.com/blusewang/pg)

#### Go Version Support
![Go version](https://img.shields.io/badge/Go-1.11-brightgreen.svg)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fblusewang%2Fpg.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fblusewang%2Fpg?ref=badge_shield)

#### PostgreSQL Version Support
![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-10.5-brightgreen.svg)
![PostgreSQL version](https://img.shields.io/badge/PostgreSQL-11.4-brightgreen.svg)


## 安装

	go get github.com/blusewang/pg

## 使用
```golang

	db, err := sql.Open("pg", "pg://user:password@dbhost.yourdomain.com/database_name?application_name=app_name&sslmode=verify-full")
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query("select * from bluse where id>$1", 0)
	if err != nil {
		return err
	}
	...

```

## 文档

更多的细节及使用示例，参见： <https://godoc.org/github.com/blusewang/pg>.


## 特性

* 送入`Scan()`处理`null`值时支持传指针类型的指针！！
* 常见`Array`类型直接兼容golang的数组类型。如PG的：`integer[]`，对应golang的：`[]int64`
* 数据源格式，既支持键值对，又支持URI。书写格式遵守：[PG官方规范](https://www.postgresql.org/docs/11/libpq-connect.html#LIBPQ-CONNSTRING)。
   * URI格式，支持`pg://`前缀。
   * 其中用户名、端口、主机名，在数据源中未指定时，有默认值。用户名默认为操作系统当前用户的用户名
   * DSN配置中，`strict`项是独立于PG后端之外的。它默认为`true`。
      * 若置为`false`；在遇到`null`值时，宽容处理。例：向`Scan()`中传 `string`型的指针，得到 `""`，传 `*string`型的指针，得到 `""`！
* 积极标记并缓存所有预备语句[包括`db.Query`、`db.Exec`、`db.Prepare()`等的语句]，遇到相同的语句请求时，自动复用。**这能提高1倍的执行速度！！！**
   * 为了发挥好此功能，需要最大可能地允许数据库连接空闲。
   * 配置上推荐将`sql.SetMaxIdleConns(x)`、`sql.SetMaxOpenConns(x)`两处的x设置为相同的值！

## 协议实现
- 此驱动更适合服务于Web

| 状态   | 功能 | 备注 |
| :---- | :---- | :---- |
| <ul><li>- [x] </li></ul> | 启动 | 必备，实现：无密码，明文密码和md5密码三种认证 |
| <ul><li>- [x] </li></ul> | 简单查询 | 必备 |
| <ul><li>- [x] </li></ul> | 扩展查询 | 必备 |
| <ul><li>- [x] </li></ul> | 取消正在处理的请求 | 必备 |
| <ul><li>- [x] </li></ul> | 终止 | 必备 |
| <ul><li>- [x] </li></ul> | SSL会话加密 | 远程安全 |


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fblusewang%2Fpg.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fblusewang%2Fpg?ref=badge_large)