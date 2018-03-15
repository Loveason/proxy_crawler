package resolver

import (
	"github.com/parnurzeal/gorequest"
	"strings"
	"time"
	"github.com/astaxie/beego/logs"
)

func getContent(addr, host string) (resp gorequest.Response, content string, errs []error) {
	req := gorequest.New().Timeout(time.Duration(10) * time.Second).
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8").
		Set("Accept-Encoding", "gzip, deflate").
		Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7").
		Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36").
		Set("Host", host)
	resp, content, errs = req.Get(addr).End()
	reqTimes := 0
	for (errs != nil || strings.TrimSpace(content) == "-10") && reqTimes < MAX_RETRY_TIMES {
		time.Sleep(time.Duration(2) * time.Second)
		logs.Error("getContent error.content:%s", content)
		resp, content, errs = getContent(addr, host)
		reqTimes++
	}
	return
}
