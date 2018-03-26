package storage

import (
	"proxy_crawler/models"
	"sync"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/parnurzeal/gorequest"
	"time"
)

var (
	db *sqlx.DB
)

func InitDB(dataSource string) error {
	var err error
	db, err = sqlx.Open("mysql", dataSource)
	if err != nil {
		return err
	}

	err = db.Ping()
	return err
}

func CheckProxy(ip *models.IP) {
	if CheckIP(ip) {
		logs.Trace("good ip:%+v", ip)
		ProxyAdd(ip)
	} else {
		logs.Trace("bad ip:%+v", ip)
	}
}

func CheckIP(ip *models.IP) bool {
	pollURL := "http://httpbin.org/get"
	resp, _, errs := gorequest.New().Timeout(time.Second * 5).Proxy("http://" + ip.Url).Get(pollURL).End()
	if errs != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func CheckProxyDB() {
	var ips []*models.IP

	err := db.Select(&ips, `select * from t_proxy`)
	if err != nil {
		logs.Trace("select error:%s", err.Error())
	}
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(v *models.IP) {
			if !CheckIP(v) {
				logs.Trace("bad ip del.%+v", v)
				ProxyDel(v)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
}

func ProxyAdd(ip *models.IP) {
	if _, err := db.Exec(`insert into t_proxy (url,type,src)values(?,?,?)`, ip.Url, ip.Type, ip.Src); err != nil {
		logs.Error("proxy add error:%s", err.Error())
	}
}

func ProxyDel(ip *models.IP) {
	if _, err := db.Exec(`delete from t_proxy where url=?`, ip.Url); err != nil {
		logs.Error("proxy del error:%s", err.Error())
	}
}
