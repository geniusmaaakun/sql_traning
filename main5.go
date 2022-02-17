package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

/*
	複数クエリがあってロールバックする必要が出てくる場合があるならトランザクションを使う。

	2,
	使うメソッドは4つで名前の通です。
	・db.Begin()
	・tx.Exec()
	・tx.Rollback()
	・tx.Commit()

	main関数で書いたのでerror処理やdeferなどざっくりですが、こんな感じです。
*/

type User struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func main() {
	db, err := sql.Open("postgres", "user=genius dbname=test_db password=genius0610 sslmode=disable")
	if err != nil {
		log.Fatal("OpenError: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("PingError: ", err)
	}

	err = Insert2(db, "Jane")
	if err != nil {
		fmt.Println(err)
	}

	var u2 User
	err = db.QueryRow("SELECT * FROM test_user WHERE name = $1", "Jane").Scan(&u2.ID, &u2.Name)
	if err != nil {
		log.Println(err)
	}

	log.Println(u2)

}

func Insert2(db *sql.DB, name string) error {
	//トランザクションはこのようにラップする。間違えなくなる（main4.go)の続き
	return Transaction(db, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO test_user(name) VALUES($1)", name)
		if err != nil {
			return fmt.Errorf("Insert() got errors: %v", err)
		}
		return nil
	})
}

//トランザクションの処理
func Transaction(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			rerr := tx.Rollback()
			err = fmt.Errorf("original=%v, rerr=%v", err, rerr)
		} else {
			if rerr := tx.Commit(); rerr != nil {
				err = fmt.Errorf("original=%v, rerr=%v", err, rerr)
			}
		}
	}()

	err = txFunc(tx)
	return err
}

/*
test=# select * from test_user where name = 'Jane';
 10 | Jane
*/
