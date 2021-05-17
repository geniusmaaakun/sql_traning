package main

/*
The Connection Pool　重要
連続したステートメントを実行すると、それぞれのコネクションがオープンされる(同一のコネクションでやる必要がある場合はトランザクションを使う)
コネクション作成数に上限がないため、too many conenctions が発生する場合がある
db.SetMacOpenConns() でコネクション最大数の制限を設けることが出来る
db.SetMaxIdleConns() でコネクションプールのアイドル数を設定で、数字が高いほどコネクションプールの再利用性が上がる
アイドル時間が長くて問題が発生する場合は db.SetMaxIdleConns(0) によって解決するかもしれない
*/

/*
Suprises, Antipatterns and Limitations　重要
リソースの枯渇 次のような使い方はしない方がいいようです
・頻繁なdatabaseとのopen(), close()をした場合
・rowsの読み込みに失敗したり、rows.Close()ができなかった場合
・rowsが帰ってこないのにQuery()をした場合(Exec()を使うべき)
・プリペアドステートメントを有効に使えない場合

uint64を利用して数値が大きすぎる場合

_, err := db.Exec("INSERT INTO users(id) VALUES", math.MaxUint64) // error
コネクション状態の不一致
・トランザクションを利用しない限りは一連のステートメントが異なるコネクションで実行される場合があります
・そのため User, BEGIN, CLOSE などを使っても異なるコネクションになったり、それをプールに返した後、次のコネクションがそれを再利用したことで不具合が発生する場合があるため、コネクションに影響が出るものは直接使ってはいけません

データベース固有のシンタックス
・データベースドライバによって実装が異なる場合があるため注意が必要です

Multiple Result Sets
・単一クエリによるMultiple Result Setsのサポートはなさそう？
・そのためストアドプロシージャで複数の結果を返す方法に対応できません

ストアドプロシージャを呼び出す
・ドライバーのサポートに依存します


複数ステートメントのサポート ・database/sql では複数のステートメントをサポートしていません
	うまく動くときもあれば、変な動きをするときもある。これらは全て予想外となります
	_, err := db.Exec("DELETE FROM tbl1; DELETE FROM tbl2") // Error/unpredictable result
	・単一のコネクションにて複数Queryを利用することはできません

	これはトランザクションではないので、別のコネクションが起動して使えます

	rows, err := db.Query("select * from tbl1") // Uses connection 1
	for rows.Next() {
		err = rows.Scan(&myvariable)
		// The following line will NOT use connection 1, which is already in-use
		db.Query("select * from tbl2 where id = ?", myvariable)
	}
	これはトランザクションなので、同一コネクションで処理しようとしてbusy状態になります

	tx, err := db.Begin()
	rows, err := tx.Query("select * from tbl1") // Uses tx's connection
	for rows.Next() {
		err = rows.Scan(&myvariable)
		// ERROR! tx's connection is already busy!　for外でやる
		tx.Query("select * from tbl2 where id = ?", myvariable)
	}
	Related Reading and Resources
	関連ドキュメントの一覧
	追加: コネクションプールの解放条件
	jmoironさんのブログを参考に抜粋。ソースコード読んで確かめるのはまたそのうちします。
	http://jmoiron.net/blog/gos-database-sql/
	sql.Row()のコネクション解放は Scan() が呼ばれたとき
	sql.Rows()のコネクション解放は Next() の終了か、 Close()の時
	sql.Txのコネクション解放はCommit()かRollback()の時
*/
