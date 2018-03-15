package models

type IP struct {
	Id   int64  `db:"id"`
	Url  string `db:"url"`
	Type string `db:"type"`
	Src  string `db:"src"`
}
