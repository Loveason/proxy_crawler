package resolver

import (
	"fmt"
	"proxy_crawler/models"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
)

const (
	KUAIDAILI       = "https://www.kuaidaili.com/free/inha"
	KDL_HOST        = "www.kuaidaili.com"
	MAX_RETRY_TIMES = 3
)

func KDL(ipChan chan<- *models.IP) {

	logs.Trace("kuaidaili")
	resp, _, errs := getContent(KUAIDAILI, KDL_HOST)
	if errs != nil {
		logs.Error("request [%s] errs:%+v", KUAIDAILI, errs)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error("read doc err:%+v", err)
		return
	}

	maxPageSize := getMaxPageSize(doc)
	tempIps := resolveByDoc(doc)
	for _, v := range tempIps {
		logs.Trace("kuaidaili ip:%+v", v)
		ipChan <- v
	}

	for pageIndex := maxPageSize; pageIndex > 1; pageIndex-- {
		addr := fmt.Sprintf("%s/%d", KUAIDAILI, pageIndex)
		logs.Trace("pageIndex:%d,addr:%s", pageIndex, addr)
		tempIps = resolveByUrl(addr)
		if len(tempIps) > 0 {
			for _, v := range tempIps {
				ipChan <- v
			}
		}
		time.Sleep(5 * time.Second)
	}

	return
}

func resolveByUrl(addr string) (result []*models.IP) {
	resp, _, errs := getContent(addr, KDL_HOST)
	if errs != nil {
		logs.Error("request [%s] errs:%+v", addr, errs)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		logs.Error("read doc err:%+v", err)
		return
	}

	result = resolveByDoc(doc)
	return
}

func resolveByDoc(doc *goquery.Document) (result []*models.IP) {
	doc.Find("#list > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		ipNode := s.Find("td[data-title='IP']")
		portNode := s.Find("td[data-title='PORT']")
		typeNode := s.Find("td[data-title='类型']")
		item := &models.IP{
			Url:  fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
			Type: typeNode.Text(),
			Src:  KDL_HOST,
		}
		result = append(result, item)
	})
	return
}

func getMaxPageSize(doc *goquery.Document) int {
	node := doc.Find("#listnav > ul > li > a").Last()
	strMaxPageSize := node.Text()
	maxPageSize, err := strconv.Atoi(strMaxPageSize)
	if err != nil {
		logs.Error("convert maxPageSize to int err:%s", err.Error())
		maxPageSize = 1
	}
	logs.Trace("maxPageSize:%s", strMaxPageSize)
	return maxPageSize
}
