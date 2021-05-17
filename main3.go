package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

/*
SQLに引数を利用する場合はdb.Prepare()かdb.Query()で渡すこと
プレースホルダで記載することでSQLインジェクション対策にもなります。もし検索機能などで外部からVALUEを受け取る場合はfmt.Sprintf()では作らないようにする。

プレースホルダ 【placeholder】
プレースホルダとは、実際の内容を後から挿入するために、とりあえず仮に確保した場所のこと。
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

	//SQLインジェクション対策で
	//これではなく。この続きで　%s部に ;で強制的に閉じられて、不当なクエリを実行される可能性がある。
	cmd := fmt.Sprintf("SELECT * FROM test_user WHERE name = %s", "'tom'")
	fmt.Println(cmd)
	rows, err := db.Query(cmd)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//どちらかでやること
	/*
		prep, err := db.Prepare("SELECT * FROM test_user WHERE name = $1")
		defer prep.Close()
		rows, err := prep.Query("tom")
		defer rows.Close()


		rows, err := db.Query("SELECT * FROM test_user WHERE name = $1", "tom")
		defer rows.Close()
	*/

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

	//シングルセレクト
	var u2 User
	err = db.QueryRow("SELECT * FROM test_user WHERE id = $1", 1).Scan(&u2.ID, &u2.Name)
	if err != nil {
		log.Println(err)
	}
	/* ------- */
	var u3 User
	prep, err := db.Prepare("SELECT * FROM test_user WHERE id = $1")
	defer prep.Close()
	err = prep.QueryRow(1).Scan(&u3.ID, &u3.Name)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(u2, u3)
}
