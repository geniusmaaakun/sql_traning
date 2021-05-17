package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

/*
DBの各操作には必ずerrorを返すため、無視しないこと
rows.Next()は問題が起きた場合に自動でrows.Close()する
func (rs *Rows) Next() bool {
	var doClose, ok bool
	withLock(rs.closemu.RLocker(), func() {
		doClose, ok = rs.nextLocked()
	})
	if doClose {
		rs.Close()
	}
	return ok
}
rows.Err()で、rows.Next()によるループ中のエラーをキャッチする
func (rs *Rows) nextLocked() (doClose, ok bool) {
	---snip---
	rs.lasterr = rs.rowsi.Next(rs.lastcols) // <---ここでclose判定 + Err()で取り出す要素にerrを入れる
	if rs.lasterr != nil {
		if rs.lasterr != io.EOF {
			return true, false
		}
		nextResultSet, ok := rs.rowsi.(driver.RowsNextResultSet)
		if !ok {
			return true, false
		}
		if !nextResultSet.HasNextResultSet() {
			doClose = true
		}
		return doClose, false
	}
	return false, true
}
rows.Close()もerrを返すが、それによって何をすべきがいいかわからないのでエラー処理の唯一の例外（強いていうならログ出力ぐらい）

検索結果がゼロの場合、QueryRow()にはエラー処理(sql.ErrNoRows)が必要。Query()は大丈夫。
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

	var u User
	err = db.QueryRow("SELECT * FROM test_user where name = $1", "John").Scan(&u.ID, &u.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("I got err but not problem: %s", err)
		} else {
			log.Fatal(err)
		}
	}

	//?
	/*
		DBのエラー判定もcastを利用することで可能な場合がある
		errの文中に含まれるかという判定も一つですが、driverによってはtypeをcastして判定する方法も用意されているようです。
		2020/08/18 15:16:15 pq: password authentication failed for user "root" <--- こういうの
		コネクションエラーの処理は特に実装の必要はない
		database/sql がコネクションプーリングの一部として最大10回のリトライ機構を備えているようです
	*/
	err = db.QueryRow("SELECT * FROM test_user where name = $1", "bob").Scan(&u.ID, &u.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("I got err but not problem: %s", err)
		} else {
			log.Fatal(err)
		}
	}
}
