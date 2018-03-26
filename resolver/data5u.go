package resolver

import (
	"proxy_crawler/models"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

type Data5u struct {
	BaseResover
}

func NewData5u() *Data5u {
	dl := &Data5u{}
	dl.r = dl
	dl.addrList = []string{
		"http://www.data5u.com/free/gngn/index.shtml",
		"http://www.data5u.com/free/gnpt/index.shtml",
	}
	dl.host = "www.data5u.com"
	return dl
}

func (dl *Data5u) GetMaxPageSize(doc *goquery.Document) (size int, err error) {
	size = 1
	return
}
func (dl *Data5u) GetPageSizeUrl(addr string, pageIndex int) string {
	return ""
}

func (dl *Data5u) ResolveContent(doc *goquery.Document) (result []*models.IP) {
	doc.Find("body > div:nth-child(7) > ul > li:nth-child(2) > ul.l2").Each(func(i int, s *goquery.Selection) {
		ipNode := s.Find("span:nth-child(1)>li")
		portNode := s.Find("span:nth-child(2)>li")
		typeNode := s.Find("span:nth-child(4)>li")
		item := &models.IP{
			Url:  fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
			Type: typeNode.Text(),
			Src:  dl.host,
		}
		result = append(result, item)
	})
	return
}