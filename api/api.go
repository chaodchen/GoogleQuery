package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"fish666/tool"
)

func init() {

}

func NewParameter(para UIParameter) DefaultParameter {
	if para.Time != "all" {
		para.Time = string(para.Time[0])
	}
	return DefaultParameter{
		hl:            "zh-CN",
		as_q:          para.Word,
		as_epq:        "",
		as_oq:         "",
		as_eq:         "",
		as_nlo:        "",
		as_nhi:        "",
		lr:            "",
		cr:            "",
		as_qdr:        para.Time,
		as_sitesearch: para.Web,
		as_occt:       para.Type,
		safe:          "image",
		as_filetype:   "",
		tbs:           "",
	}
}

func GetSearchRet(p UIParameter, back func(string, error)) {
	request_sleep, err := strconv.Atoi(tool.ReadIni("http", "request_sleep"))
	if err != nil {
		request_sleep = 0
	}
	fmt.Printf("[*] request_sleep: %d", request_sleep)
	// time.Sleep(time.Duration(request_sleep) * time.Second)

	searchApi := "https://www.google.com/search?"
	client := &http.Client{Timeout: 15 * time.Second}
	if p.Proxy != "" {
		fmt.Printf("[*] proxy: %s\n", p.Proxy)
		proxyURL, err := url.Parse(p.Proxy)
		if err != nil {
			fmt.Println("ProxyUrl解析失败")
			back("", err)
			return
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	queryParams := url.Values{}
	np := NewParameter(p)

	queryParams.Set("hl", np.hl)
	queryParams.Set("as_q", np.as_q)
	queryParams.Set("as_epq", np.as_epq)
	queryParams.Set("as_oq", np.as_oq)
	queryParams.Set("as_eq", np.as_eq)
	queryParams.Set("as_nlo", np.as_nlo)
	queryParams.Set("as_nhi", np.as_nhi)
	queryParams.Set("lr", np.lr)
	queryParams.Set("as_qdr", np.as_qdr)
	queryParams.Set("as_sitesearch", np.as_sitesearch)
	queryParams.Set("as_occt", np.as_occt)
	queryParams.Set("safe", np.safe)
	queryParams.Set("as_filetype", np.as_filetype)
	queryParams.Set("tbs", np.tbs)

	tagerLink := searchApi + queryParams.Encode()
	fmt.Println(tagerLink)
	req, _ := http.NewRequest("GET", tagerLink, nil)

	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh-Hans;q=0.9")
	// req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", tool.ReadIni("http", "user_agent"))
	fmt.Println("[*] Do Requests")
	resp, err := client.Do(req)
	fmt.Println("[*] Requests end.")
	if err != nil {
		fmt.Printf("[*] Requests error: %s\n", err.Error())
		back("", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("[*] http.Status not equal to ok.")
		back("", nil)
		return
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		back("", err)
		fmt.Println("[*] io.Readall error.")
		return
	}

	// 正则取值
	re := regexp.MustCompile(`<div id="result-stats">找到约 ([0-9,]+) 条结果<nobr> （用时 ([0-9.]+) 秒）`)
	match := re.FindStringSubmatch(string(body))
	if len(match) > 2 {
		result := match[1]
		duration := match[2]
		fmt.Printf("[*] site: %s, 结果数量: %s, 用时%s\n", p.Web, result, duration)
		back(result, nil)
		return
	} else {
		fmt.Println("[!] not found search result.")
		file, _ := os.Create("failed.html")
		defer file.Close()
		file.WriteString(string(body))
		back("", nil)
		return
	}
}
