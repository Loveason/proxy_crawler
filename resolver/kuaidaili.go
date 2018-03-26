package resolver

import (
	"fmt"
	"proxy_crawler/models"
	"strconv"
	"github.com/PuerkitoBio/goquery"
)

type KuaiDaili struct {
	BaseResover
}

func NewKuaiDaili() *KuaiDaili {
	dl := &KuaiDaili{}
	dl.r = dl
	dl.addrList = []string{
		"https://www.kuaidaili.com/free/inha",
	}
	dl.host = "www.kuaidaili.com"
	return dl
}

func (dl *KuaiDaili) GetMaxPageSize(doc *goquery.Document) (size int, err error) {
	maxPageNode := doc.Find("#listnav > ul > li > a").Last()
	size, err = strconv.Atoi(maxPageNode.Text())
	return
}
func (dl *KuaiDaili) GetPageSizeUrl(addr string, pageIndex int) string {
	return fmt.Sprintf("%s/%d", addr, pageIndex)
}

func (dl *KuaiDaili) ResolveContent(doc *goquery.Document) (result []*models.IP) {
	doc.Find("#list > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		ipNode := s.Find("td[data-title='IP']")
		portNode := s.Find("td[data-title='PORT']")
		typeNode := s.Find("td[data-title='类型']")
		item := &models.IP{
			Url:  fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
			Type: typeNode.Text(),
			Src:  dl.host,
		}
		result = append(result, item)
	})
	return
}
//
//
//func KDL(ipChan chan<- *models.IP) {
//
//	logs.Trace("kuaidaili")
//
//	resp, _, errs := getContent(KUAIDAILI, KDL_HOST)
//	if errs != nil {
//		logs.Error("request [%s] errs:%+v", KUAIDAILI, errs)
//		return
//	}
//
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		logs.Error("read doc err:%+v", err)
//		return
//	}
//
//	maxPageSize := getMaxPageSize(doc)
//	tempIps := resolveByDoc(doc)
//	for _, v := range tempIps {
//		logs.Trace("kuaidaili ip:%+v", v)
//		ipChan <- v
//	}
//
//	for pageIndex := maxPageSize; pageIndex > 1; pageIndex-- {
//		addr := fmt.Sprintf("%s/%d", KUAIDAILI, pageIndex)
//		logs.Trace("pageIndex:%d,addr:%s", pageIndex, addr)
//		tempIps = resolveByUrl(addr)
//		if len(tempIps) > 0 {
//			for _, v := range tempIps {
//				ipChan <- v
//			}
//		}
//		time.Sleep(5 * time.Second)
//	}
//
//	return
//}
//
//func resolveByUrl(addr string) (result []*models.IP) {
//	_, content, errs := getContent(addr, KDL_HOST)
//	if errs != nil {
//		logs.Error("request [%s] errs:%+v", addr, errs)
//		return
//	}
//
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
//
//	if err != nil {
//		logs.Error("read doc err:%+v", err)
//		return
//	}
//
//	result = resolveByDoc(doc)
//	return
//}
//
//func resolveByDoc(doc *goquery.Document) (result []*models.IP) {
//	doc.Find("#list > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
//		ipNode := s.Find("td[data-title='IP']")
//		portNode := s.Find("td[data-title='PORT']")
//		typeNode := s.Find("td[data-title='类型']")
//		item := &models.IP{
//			Url:  fmt.Sprintf("%s:%s", ipNode.Text(), portNode.Text()),
//			Type: typeNode.Text(),
//			Src:  KDL_HOST,
//		}
//		result = append(result, item)
//	})
//	return
//}
//
//func getMaxPageSize(doc *goquery.Document) int {
//	node := doc.Find("#listnav > ul > li > a").Last()
//	strMaxPageSize := node.Text()
//	maxPageSize, err := strconv.Atoi(strMaxPageSize)
//	if err != nil {
//		logs.Error("convert maxPageSize to int err:%s", err.Error())
//		maxPageSize = 1
//	}
//	logs.Trace("maxPageSize:%s", strMaxPageSize)
//	return maxPageSize
//}
