package resolver

import (
	"github.com/astaxie/beego/logs"
	"proxy_crawler/models"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"io/ioutil"
	"io"
	"github.com/parnurzeal/gorequest"
	"time"
	"compress/gzip"
	"compress/zlib"
)

type IResove interface {
	RequestContent(addr string) (content string, err error)
	GetMaxPageSize(doc *goquery.Document) (size int, err error)
	GetPageSizeUrl(addr string, pageIndex int) string
	ResolveContent(doc *goquery.Document) (result []*models.IP)
}

type BaseResover struct {
	r        IResove
	addrList []string
	host     string
}

func (b *BaseResover) Resolve(ipChan chan<- *models.IP) {
	defer func() {
		logs.Trace("%s解析完成", b.host)
	}()

	for _, addr := range b.addrList {
		content, err := b.r.RequestContent(addr)
		if err != nil {
			logs.Error(err)
		}
		if strings.TrimSpace(content) == "block" {
			logs.Error("%s被禁", b.host)
			return
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return
		}
		ips := b.r.ResolveContent(doc)
		if len(ips) == 0 {
			logs.Error("resolve error.addr:%s", addr)
		} else {
			logs.Trace("resolve good.addr:%s", addr)
		}
		logs.Trace("ips:%+v", ips)
		for _, ip := range ips {
			ipChan <- ip
		}

		pageSize, err := b.r.GetMaxPageSize(doc)
		if err != nil {
			logs.Error("getMaxPageSize error:%s", err.Error())
		}
		logs.Trace("pageSize:%d", pageSize)
		if pageSize > 1 {
			for i := pageSize; i > 2; i-- {
				pagedAddr := b.r.GetPageSizeUrl(addr, i)
				content, err = b.r.RequestContent(pagedAddr)
				if err != nil {
					logs.Error(err)
				}
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
				if err != nil {
					return
				}
				ips = b.r.ResolveContent(doc)
				if len(ips) == 0 {
					logs.Error("resolve error.addr:%s", pagedAddr)
				} else {
					logs.Trace("resolve good.addr:%s", pagedAddr)
				}
				for _, ip := range ips {
					ipChan <- ip
				}
				time.Sleep(time.Second * 2)
			}
		}
	}
}

func (b *BaseResover) RequestContent(addr string) (content string, err error) {
	var (
		reader io.ReadCloser
	)
	resp, _, errs := gorequest.New().Timeout(time.Duration(10) * time.Second).Get(addr).
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8").
		Set("Accept-Encoding", "gzip, deflate").
		Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7").
		Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36").
		End()
	if errs != nil {
		err = errs[0]
		return
	}

	enc := resp.Header.Get("Content-Encoding")

	switch enc {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
		defer reader.Close()
	case "deflate":
		reader, err = zlib.NewReader(resp.Body)
		if err != nil {
			return
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}
	bt, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	content = string(bt)

	return
}
