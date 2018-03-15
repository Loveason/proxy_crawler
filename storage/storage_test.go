package storage

import (
	"testing"

	"proxy_crawler/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func TestSqlxSelect(t *testing.T) {
	var err error
	db, err = sqlx.Open("mysql", `root:wbs^1XTRz6a0UrQb@tcp(47.97.110.163:3306)/proxyDB?charset=utf8`)
	if err != nil {
		t.Error("open err:%s", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		t.Error("ping err:%s", err.Error())
		return
	}

	ips := []*models.IP{}
	//var urls []string

	err = db.Select(&ips, `select * from t_proxy`)
	if err != nil {
		t.Error("select err:", err.Error())
		return
	}

	for _, ip := range ips {
		t.Log("ip:", ip)
	}

}
