package resolver

import (
	"proxy_crawler/models"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"fmt"
)

type IP66 struct {
	BaseResover
}

func NewIP66() *IP66 {
	dl := &IP66{}
	dl.r = dl
	dl.addrList = []string{
		"http://www.66ip.cn/",
	}
	dl.host = "www.66ip.cn"
	return dl
}

func (dl *IP66) GetMaxPageSize(doc *goquery.Document) (size int, err error) {
	maxPageNode := doc.Find("#PageList > a:nth-last-child(2)")
	size, err = strconv.Atoi(maxPageNode.Text())
	return
}
func (dl *IP66) GetPageSizeUrl(addr string, pageIndex int) string {
	return fmt.Sprintf("%s/%d.html", addr, pageIndex)
}

func (dl *IP66) ResolveContent(doc *goquery.Document) (result []*models.IP) {
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

//func Ip66(ipChan chan<- *models.IP) {
//	logs.Trace("66ip")
//
//	_, content, errs := getContent(IP66, IP66_HOST)
//	if errs != nil {
//		logs.Error("get err.url:%s,errs:%+v", IP66, errs)
//		return
//	}
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
//
//	if err != nil {
//		logs.Error("read doc err:%+v", err)
//		return
//	}
//
//	temp := ip66Resolve(doc)
//	for _, v := range temp {
//		logs.Trace("66ip ip:%+v", v)
//		ipChan <- v
//	}
//
//	maxPageNode := doc.Find("#PageList > a:nth-last-child(2)")
//	maxPageSize, err := strconv.Atoi(maxPageNode.Text())
//
//	for i := maxPageSize; i > 1; i-- {
//		addr := fmt.Sprintf("%s%d.html", IP66, i)
//		_, content, errs = getContent(addr, IP66_HOST)
//
//		if errs != nil {
//			logs.Error("get err.url:%s,errs:%+v", IP66, errs)
//			continue
//		}
//		doc, err = goquery.NewDocumentFromReader(strings.NewReader(content))
//
//		if err != nil {
//			logs.Error("read doc err:%+v", err)
//			continue
//		}
//		temp = ip66Resolve(doc)
//		for _, v := range temp {
//			ipChan <- v
//		}
//	}
//	return
//}

//func ip66Resolve(doc *goquery.Document) (result []*models.IP) {
//
//	doc.Find("#main > div > div:nth-child(1) > table > tbody > tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
//		ipNode := s.Find("td:nth-child(1)")
//		portNode := s.Find("td:nth-child(2)")
//		item := &models.IP{
//			Url: fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
//			Src: IP66_HOST,
//		}
//		result = append(result, item)
//	})
//
//	return
//}
