package resolver

import (
	"proxy_crawler/models"
	"time"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"fmt"
)

const (
	XCDL_HOST = "www.xicidaili.com"
)

func Xicidaili(ipChan chan<- *models.IP) {

	logs.Trace("xicidaili")
	addrs := []string{
		"http://www.xicidaili.com/nn/",
		"http://www.xicidaili.com/nt/",
		"http://www.xicidaili.com/wn/",
		"http://www.xicidaili.com/wt/",
	}

	for _, addr := range addrs {

		resp, content, errs := getContent(addr, XCDL_HOST)
		if errs != nil {
			logs.Error("get err.url:%s,errs:%+v", addr, errs)
			continue
		}
		logs.Trace(content)
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			logs.Error("read doc err:%+v", err)
			continue
		}
		temp := xiciResolve(doc)
		for _, v := range temp {
			logs.Trace("xicidaili ip:%+v", v)
			ipChan <- v
		}

		pageNodes := doc.Find("#body > div.pagination > a").Nodes
		logs.Trace("pageNodes length:%d", len(pageNodes))
		for _, v := range pageNodes {
			logs.Trace("node:%s", v.Data)
		}
		maxPageNode := doc.Find("#body > div.pagination > a:nth-last-child(2)")

		maxPageSize, err := strconv.Atoi(maxPageNode.Text())
		if err != nil {
			logs.Error("atoi err:%s,maxPageNodeText:%s,addr:%s", err.Error(), maxPageNode.Text(), addr)
			continue
		}

		for i := maxPageSize; i > 1; i-- {
			addr = fmt.Sprintf("%s%d", addr, i)
			resp, content, errs = getContent(addr, XCDL_HOST)
			logs.Trace(content)
			if errs != nil {
				logs.Error("get err.url:%s,errs:%+v", addr, errs)
				continue
			}

			doc, err = goquery.NewDocumentFromResponse(resp)
			if err != nil {
				logs.Error("read doc err:%+v", err)
				continue
			}
			temp = xiciResolve(doc)
			for _, v := range temp {
				ipChan <- v
			}
			time.Sleep(2 * time.Second)
		}
	}

	return
}

func xiciResolve(doc *goquery.Document) (result []*models.IP) {

	doc.Find("#ip_list > tbody > tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
		ipNode := s.Find("td:nth-child(2)")
		portNode := s.Find("td:nth-child(3)")
		typeNode := s.Find("td:nth-child(6)")

		item := &models.IP{
			Url:  fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
			Type: typeNode.Text(),
			Src:  XCDL_HOST,
		}
		result = append(result, item)
	})

	return
}
