package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

/*
Retrieving Results Sets
大きくQuery()とExec()の2種類あり、次の用途で使い分けること
・Query: 複数の検索結果(rows)を取得する場合(SELECT)。 一行(row)の場合はQueryRow()。
・Exec: 検索結果を取得しない場合(CREATE, INSERT, UPDATE, DELETE etc)

SELECTを試すために、こんな感じで適当にデータをDBに入力します。正しくは次の Modifying Data and Using Transactions を参照ください。

rows.Close()は rows.Err() のあとで、defer を for外 (memory leakする)で必ず利用
rows.Scan()はpointerで渡した型に取得したdataをマッピング
上記の例ではScan(&u.ID, &u.Name)と構造体の中身を一つずつ渡していますが、数が増えると正直手間です。これを構造体渡すだけ(&u)でいいようにするには package database/sqlx が使えます。

SQLに引数を利用する場合はdb.Prepare()かdb.Query()で渡すこと
プレースホルダで記載することでSQLインジェクション対策にもなります。もし検索機能などで外部からVALUEを受け取る場合はfmt.Sprintf()では作らないようにする。

プレースホルダ 【placeholder】
プレースホルダとは、実際の内容を後から挿入するために、とりあえず仮に確保した場所のこと。
*/

type User struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func Insert(db *sql.DB) error {
	//db作成
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS test_user (
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(10))`)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO test_user (name) VALUES ('tom')")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := sql.Open("postgres", "user=genius dbname=test_db password=genius0610 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//db作成、データ挿入
	err = Insert(db)
	if err != nil {
		log.Println(err)
	}

	rows, err := db.Query("SELECT * FROM test_user")
	if err != nil {
		log.Println(err)
	}

	//rows.Close()は rows.Err() のあとで、defer を for外 (memory leakする)で必ず利用
	defer rows.Close()

	var u User
	//rows.Scan()はpointerで渡した型に取得したdataをマッピング
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s, Name: %s\n", u.ID, u.Name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
