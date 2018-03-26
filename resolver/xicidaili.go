package resolver

import (
	"proxy_crawler/models"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"fmt"
)

type XiciDaili struct {
	BaseResover
}

func NewXiciDaili() *XiciDaili {
	dl := &XiciDaili{}
	dl.r = dl
	dl.addrList = []string{
		"http://www.xicidaili.com/nn/",
		"http://www.xicidaili.com/nt/",
		"http://www.xicidaili.com/wn/",
		"http://www.xicidaili.com/wt/",
	}
	dl.host = "www.xicidaili.com"
	return dl
}

func (dl *XiciDaili) GetMaxPageSize(doc *goquery.Document) (size int, err error) {
	maxPageNode := doc.Find("#body > div.pagination > a:nth-last-child(2)")
	size, err = strconv.Atoi(maxPageNode.Text())
	return
}
func (dl *XiciDaili) GetPageSizeUrl(addr string, pageIndex int) string {
	return fmt.Sprintf("%s/%d", addr, pageIndex)
}

func (dl *XiciDaili) ResolveContent(doc *goquery.Document) (result []*models.IP) {
	doc.Find("#ip_list > tbody > tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
		ipNode := s.Find("td:nth-child(2)")
		portNode := s.Find("td:nth-child(3)")
		typeNode := s.Find("td:nth-child(6)")

		item := &models.IP{
			Url:  fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
			Type: typeNode.Text(),
			Src:  dl.host,
		}
		result = append(result, item)
	})

	return
}
