package main

/*
Using Prepared Statements
Goではプリペアドステートメントをdb.Query(), db.Prepare(), tx.Prepare()などで利用で
プリペアするとDB poolのconnectionに紐づいたStmtを返すため、それを利用する
仮に紐づいたコネクションが使えなかったら新たに内部で取得する
ただ言い換えるとビジー状態が続いてしまうと、プリペアドステートメントを大量に作るので気づいたら上限に達する可能性がある
プリペアドステートメントを利用したくない場合はdb.Query(fmt.Sprintf(str))で渡す
*/