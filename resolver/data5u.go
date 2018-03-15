package resolver

import (
	"proxy_crawler/models"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

const (
	DATA5U_HOST = "www.data5u.com"
)

func Data5u(ipChan chan<- *models.IP) {

	logs.Trace("data5u")

	addrs := []string{
		"http://www.data5u.com/free/gngn/index.shtml",
		"http://www.data5u.com/free/gnpt/index.shtml",
	}
	for _, addr := range addrs {
		resp, _, errs := getContent(addr, DATA5U_HOST)
		if errs != nil {
			logs.Error("get err.url:%s,errs:%+v", addr, errs)
			continue
		}

		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			logs.Error("read doc err:%+v", err)
			continue
		}

		doc.Find("body > div:nth-child(7) > ul > li:nth-child(2) > ul.l2").Each(func(i int, s *goquery.Selection) {
			ipNode := s.Find("span:nth-child(1)>li")
			portNode := s.Find("span:nth-child(2)>li")
			typeNode := s.Find("span:nth-child(4)>li")
			item := &models.IP{
				Url:  fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
				Type: typeNode.Text(),
				Src:  DATA5U_HOST,
			}
			ipChan <- item
		})
	}

	return
}
