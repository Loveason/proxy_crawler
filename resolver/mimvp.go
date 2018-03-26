package resolver

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"fmt"
	"proxy_crawler/models"
)

type MiMvp struct {
	BaseResover
}

func NewMiMvp() *MiMvp {
	dl := &MiMvp{}
	dl.r = dl
	dl.addrList = []string{
		"https://proxy.mimvp.com/free.php?proxy=in_tp",
		"https://proxy.mimvp.com/free.php?proxy=in_hp",
	}
	dl.host = "proxy.mimvp.com"
	return dl
}

func (dl *MiMvp) GetMaxPageSize(doc *goquery.Document) (size int, err error) {
	maxPageNode := doc.Find("#listnav > ul > li:nth-last-child(2) > a")
	size, err = strconv.Atoi(maxPageNode.Text())
	return
}
func (dl *MiMvp) GetPageSizeUrl(addr string, pageIndex int) string {
	return fmt.Sprintf("%s&sort=&page=%d", addr, pageIndex)
}

func (dl *MiMvp) ResolveContent(doc *goquery.Document) (result []*models.IP) {
	doc.Find("#main > div > div:nth-child(1) > table > tbody > tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
		ipNode := s.Find("td:nth-child(1)")
		portNode := s.Find("td:nth-child(2)")
		item := &models.IP{
			Url: fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
			Src: dl.host,
		}
		result = append(result, item)
	})

	return
}


