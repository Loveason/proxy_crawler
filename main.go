package main

import (
	"fmt"
	"os"
	"proxy_crawler/resolver"

	"proxy_crawler/models"
	"proxy_crawler/storage"
	"sync"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dataSource string
	cfg        config.Configer
)

func init() {
	var err error
	logs.SetLogger(logs.AdapterConsole)
	err = logs.SetLogger(logs.AdapterFile, `{"filename":"logs/crawler.log","maxdays":3}`)
	if err != nil {
		fmt.Printf("SetLogger Error:%v \n", err)
		os.Exit(1)
	}

	//配置
	cfg, err = config.NewConfig("ini", "my.cfg")
	if err != nil {
		logs.Error("读取配置文件异常.")
		os.Exit(1)
	}

	dataSource = cfg.String("mysql::dataSource")
	storage.InitDB(dataSource)
}

func main() {
	ipChan := make(chan *models.IP, 2000)
	go func() {
		for {
			storage.CheckProxyDB()
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()

	for i := 0; i < 50; i++ {
		go func() {
			for {
				storage.CheckProxy(<-ipChan)
			}
		}()
	}

	for {
		if len(ipChan) < 100 {
			run(ipChan)
			time.Sleep(10 * time.Minute)
		}
	}
}

func run(ipChan chan<- *models.IP) {
	var wg sync.WaitGroup
	funs := []func(chan<- *models.IP){
		resolver.NewXiciDaili().Resolve,
		resolver.NewIP66().Resolve,
		resolver.NewData5u().Resolve,
		resolver.NewKuaiDaili().Resolve,
	}

	for _, f := range funs {
		wg.Add(1)
		go func() {
			f(ipChan)
			wg.Done()
		}()
	}
	wg.Wait()
	logs.Trace("All resolver finished.")
}
