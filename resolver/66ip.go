package resolver

import (
	"proxy_crawler/models"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"fmt"
)

const (
	IP66      = "http://www.66ip.cn/"
	IP66_HOST = "www.66ip.cn"
)

func Ip66(ipChan chan<- *models.IP) {
	logs.Trace("66ip")
	resp, _, errs := getContent(IP66, IP66_HOST)
	if errs != nil {
		logs.Error("get err.url:%s,errs:%+v", IP66, errs)
		return
	}
	doc, err := goquery.NewDocumentFromResponse(resp)

	if err != nil {
		logs.Error("read doc err:%+v", err)
		return
	}

	temp := ip66Resolve(doc)
	for _, v := range temp {
		logs.Trace("66ip ip:%+v", v)
		ipChan <- v
	}

	maxPageNode := doc.Find("#PageList > a:nth-last-child(2)")
	maxPageSize, err := strconv.Atoi(maxPageNode.Text())

	for i := maxPageSize; i > 1; i-- {
		addr := fmt.Sprintf("%s%d.html", IP66, i)
		resp, _, errs = getContent(addr, IP66_HOST)

		if errs != nil {
			logs.Error("get err.url:%s,errs:%+v", IP66, errs)
			continue
		}
		doc, err = goquery.NewDocumentFromResponse(resp)

		if err != nil {
			logs.Error("read doc err:%+v", err)
			continue
		}
		temp = ip66Resolve(doc)
		for _, v := range temp {
			ipChan <- v
		}
	}
	return
}

func ip66Resolve(doc *goquery.Document) (result []*models.IP) {

	doc.Find("#main > div > div:nth-child(1) > table > tbody > tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
		ipNode := s.Find("td:nth-child(1)")
		portNode := s.Find("td:nth-child(2)")
		item := &models.IP{
			Url: fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
			Src: IP66_HOST,
		}
		result = append(result, item)
	})

	return
}
