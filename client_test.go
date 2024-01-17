package pg

import (
	"context"
	"database/sql"
	"github.com/blusewang/pg/v2/internal/client"
	"io"
	"log"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := client.NewClient()
	dsn, err := client.ParseDSN("pg://postgres:bluse.123@r/core?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Connect(context.Background(), dsn); err != nil {
		t.Fatal(err)
	}
	if err = c.AutoSSL(); err != nil {
		t.Fatal(err)
	}
	if err = c.Startup(); err != nil {
		t.Log(err.(Error))
		t.Fatal(err)
	}

	pr, err := c.Parse("x", "select * from pg_sleep(1,1)")
	if err != nil {
		t.Log(err)
	}
	t.Log(pr.Rows)

}

func TestSql(t *testing.T) {
	sql.Register("pq", Driver{})
	ctx := context.Background()
	db, err := sql.Open("pq", "pg://postgres:bluse.123@r:5432/core?sslmode=disable&application_name=development_w")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)

	rows, err := db.QueryContext(ctx, "select * from pg_sleep(2,2)")
	if err != nil {
		t.Log(err)
	}
	t.Log(rows)

	//for i := 0; i < 8; i++ {
	//	queryRow(d, ctx, t)
	//	time.Sleep(time.Second * 4)
	//}

	time.Sleep(time.Second)
}

func query(d *sql.DB, ctx context.Context, t *testing.T) {
	var v time.Time
	rows, err := d.QueryContext(ctx, "select CURRENT_TIME::time without time zone from shops where sid = any($1)", []int{2129, 206322})
	if err != nil {
		t.Log(err)
		return
	}
	defer rows.Close()
	t.Log(rows.Columns())
	t.Log(rows.ColumnTypes())
	for rows.Next() {
		t.Log(rows.Scan(&v))
		t.Log(v)
	}
}

func queryRow(d *sql.DB, ctx context.Context, t *testing.T) {
	t.Log("begin>>>>>>>>>>>")
	var v time.Time
	err := d.QueryRowContext(ctx, "select CURRENT_TIME::time without time zone from shops where sid = any($1)", []int{2129, 206322}).Scan(&v)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(v)
}

func TestNewListener(t *testing.T) {
	l, err := NewListener(context.Background(), "pg://developer:dev.123@mywsy.cn:5432/core?application_name=listener")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(l.Listen("public_bills"))
	go func() {
		time.Sleep(time.Minute)
		_ = l.Terminate()
	}()
	for {
		pid, channel, message, err := l.GetNotification()
		if err == io.EOF {
			t.Fatal(err)
		} else if err != nil {
			t.Fatal(err)
		}
		log.Println(pid, channel, message)
	}
}

func TestSql2(t *testing.T) {
	sql.Register("pq", Driver{})
	db, err := sql.Open("pq", "pg://postgres:bluse.123@r:5432/core?sslmode=disable&application_name=development_w")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := db.Query("select shop_name from shops where sid=any($1) xx", []int{1})
	ts, _ := rows.ColumnTypes()
	for i, columnType := range ts {
		t.Log(i, columnType.Name(), columnType.DatabaseTypeName(), columnType.ScanType())
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rows.Next())
}

func loadDb(t *testing.T) (db *sql.DB) {
	sql.Register("pq", Driver{})
	db, err := sql.Open("pq", "pg://postgres:bluse.123@r:5432/core?sslmode=disable&application_name=development_w")
	if err != nil {
		t.Fatal(err)
	}
	return
}
func TestNull(t *testing.T) {
	db := loadDb(t)
	rows, err := db.Query("select delete_at from shops where sid=2129")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var ts *time.Time
		t.Log(rows.Scan(&ts))
		t.Log(ts)
	}
	t.Log(rows.Close())
}
