package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

/*
Modifying Data and Using Transactions
Insertを使う場合はExecを使用するが、InsertされたIDを確認するにはMySQLとPostgresで方法が異なる
・MySQL
サンプルの直接抜粋です
stmt, err := db.Prepare("INSERT INTO test_user(name) VALUES(?)")
if err != nil {
	log.Fatal(err)
}
res, err := stmt.Exec("john")
if err != nil {
	log.Fatal(err)
}
LastID, err := res.LastInsertId()
if err != nil {
	log.Fatal(err)
}
rowCnt, err := res.RowsAffected()
if err != nil {
	log.Fatal(err)
}
log.Printf("ID = %d, affected = %d\n", LastID, rowCnt)


PostgreSQL
Postgresの場合は LanstInsertId() に非対応(see. #24)でPostgres自体の仕様です。

$ go run main.go
2020/08/18 10:30:57 LastInsertId is not supported by this driver
exit status 1
そのため次のように確認します。

var ID int
err = db.QueryRow("INSERT INTO test_user(name) VALUES($1) RETURNING id", "john").Scan(&ID)
if err != nil {
	log.Fatal(err)
}
log.Printf("ID = %d\n", ID)

/*
test=# select * from test_user where id = 6;
  6 | john
*/

type User struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func main() {
	db, err := sql.Open("postgres", "user=genius dbname=test_db password=genius0610 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	/*
		INSERTは db.Exec() じゃないの？と思うかもしれません。これは db.QueryRow() がrowの戻りを期待して待ってしまうことで Connectionを予約し続けることに起因します。
		一方で上記の場合は postgres から id を返すため、db.QueryRow() で実行しています。

		つまり
		1,rowsの戻りがない場合は db.Exec() を使う

		2,SQL statementを利用した直接の BEGIN, COMMIT は不具合を招くことがあるので使わない
	*/
	var ID int
	err = db.QueryRow("INSERT INTO test_user(name) VALUES($1) RETURNING id", "john").Scan(&ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d\n", ID)

	/*
		2,
		使うメソッドは4つで名前の通です。
		・db.Begin()
		・tx.Exec()
		・tx.Rollback()
		・tx.Commit()

		main関数で書いたのでerror処理やdeferなどざっくりですが、こんな感じです。
	*/
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	//直接書かない
	/*
		使ってはいけない理由は上記のメソッドを使わないと、一連のステートメントの処理は異なるコネクションを利用する場合があり、
		そもそもトランザクション処理が想定外な動作になるケースや、返却したコネクションを再利用したコネクションにその影響が出てしまうため、
		必ずトランザクション用のコネクション(Tx)を作って実行する必要があるようです。

		番外編として、次のようにするとtransactionの処理をWrapすると効率的(see. database/sql Tx - detecting Commit or Rollback)
		main5.go
	*/
	_, err = tx.Exec("INSERT INTO test_user (name) VALUES ($1)", "TRANS")
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(err)
		}
	} else {
		if err := tx.Commit(); err != nil {
			log.Println(err)
		}
	}

	var u2 User
	err = db.QueryRow("SELECT * FROM test_user WHERE name = $1", "TRANS").Scan(&u2.ID, &u2.Name)
	if err != nil {
		log.Println(err)
	}

	log.Println(u2)
}
