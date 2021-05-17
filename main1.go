package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

/*
GoによりDBへのアクセスはsql.Open()を使いますが、これはDBへアクセスできるオブジェクトを返すだけなので、
実際の接続テストはsql.Ping()を使用することで確認できます(ただしdriverによってping()の実装は変わるようです)
*/
func main() {
	db, err := sql.Open("postgres", "user=genius dbname=test_db password=genius0610 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("PingError: ", err)
	}
}

/*
psql, mysql

func main() {
	db, err := sql.Open("postgres", "user=root password=root host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal("OpenError: ", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("PingError: ", err)
	}
}

func OpenDB(c driver.Connector) *DB {
	ctx, cancel := context.WithCancel(context.Background())
	db := &DB{
		connector:    c,
		openerCh:     make(chan struct{}, connectionRequestQueueSize),
		lastPut:      make(map[*driverConn]string),
		connRequests: make(map[uint64]chan connRequest),
		stop:         cancel,
	}

	go db.connectionOpener(ctx)

	return db
}
で c（user=roo password=root host=localhost dbname=test sslmode=disable）で間違った情報を渡すと
db.Ping()により接続テスト
わざとユーザー名を間違えた場合、db.Ping()で検出されています

$ go run main.go
2020/08/17 22:47:24 PingError: pq: password authentication failed for user "roo"
exit status 1
dbはConnection Poolを利用するため一度Open()したものをClose()せずに使いまわすことが基本

Open(), Close()を頻繁にすると利用効率の低下、ネットワーク帯域の圧迫、TCPのTIME_WAITなどが発生するやも
*/
