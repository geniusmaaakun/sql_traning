package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Insert3(db *sql.DB) error {
	//db作成
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS test_user2 (
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(10),
		age INTEGER)`)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO test_user2 (name) VALUES ('tom')")
	if err != nil {
		return err
	}
	return nil
}

/*
Working with NULLs
NULLは使わない方がいい。
DBの観点もあるが、Goの観点からではゼロ値の概念があるため相性が悪い。
次のようにcastすることはできません
*/
type User2 struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

/*
test=# select * from test_user;
  1 | bob  |  12
  2 | tom  |  	  <--- NULL

go run main8.go
sql: expected 2 destination arguments in Scan, not 3
{  0}
sql: expected 2 destination arguments in Scan, not 3
{  0}
sql: expected 2 destination arguments in Scan, not 3
{  0}
sql: expected 2 destination arguments in Scan, not 3
{  0}
sql: expected 2 destination arguments in Scan, not 3
{  0}
sql: expected 2 destination arguments in Scan, not 3
{  0}
sql: expected 2 destination arguments in Scan, not 3
{  0}
sql: expected 2 destination arguments in Scan, not 3
{  0}
sql: expected 2 destination arguments in Scan, not 3
{  0}
*/

//DBにNullが紛れ込んでいるならsql.NullXXXでcastするが、構造が変わってしまうことに注意
type User3 struct {
	ID   string        `db:"id"`
	Name string        `db:"name"`
	Age  sql.NullInt32 `db:"age"`
}

/*
$ go run main8.go
{1 bob {12 true}}
{2 tom {0 false}}
*/

func main() {
	db, err := sql.Open("postgres", "user=genius dbname=test_db password=genius0610 sslmode=disable")
	if err != nil {
		log.Fatal("OpenError: ", err)
	}

	err = Insert3(db)
	if err != nil {
		log.Println(err)
	}

	var u User2
	rows, _ := db.Query("SELECT * FROM test_user2")
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			//if err := rows.Scan(&u.ID, &u.Name); err != nil {

			fmt.Println(err)
		}
		fmt.Println(u)
	}

	fmt.Println()

	var u2 User3
	rows2, _ := db.Query("SELECT * FROM test_user2")
	for rows2.Next() {
		if err := rows2.Scan(&u2.ID, &u2.Name, &u2.Age); err != nil {
			//if err := rows2.Scan(&u2.ID, &u2.Name); err != nil {

			fmt.Println(err)
		}
		fmt.Println(u2)
	}

	/*
		ageの出力が別の構造体に変わっているのは、次のように定義されているためです

		type NullInt32 struct {
			Int32 int32
			Valid bool // Valid is true if Int32 is not NULL
		}
		ややこしいのは type User は上記のデータ構造なので、単純にUser.AgeをIntで扱うことができません

		COALESCE()を利用してNULLを置き換える方法をとるtipsもある
		Databaseの機能でNULLの時にCOALESCE()で取得できた値に置き換える機能ですが、NULLに置き換えたい0や'‘などを用意することでNULLを出さないようにしています
	*/

	fmt.Println()

	var u3 User3
	rows3, _ := db.Query("SELECT id, name, COALESCE(age, 0) as age FROM test_user2")
	for rows3.Next() {
		if err := rows3.Scan(&u3.ID, &u3.Name, &u3.Age); err != nil {
			fmt.Println(err)
		}
		fmt.Println(u3)
	}
	/*
		$ go run main.go
		{1 bob 12}
		{2 tom 0}
	*/
	//これらを考えるとNULLと0や'‘を正しく区別したい場合には難しいかもしれません。。 ただ入力区分を作ってNULLかどうか判定する方法もあると思うので、データベースの設計で次第かもしれないですね。

}
