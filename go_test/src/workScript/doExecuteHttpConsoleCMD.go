package fileManage

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"io"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strconv"
)

func WebGetAction(url string) []byte {
	client := &http.Client{
		CheckRedirect: nil,
	}
	reqest, _ := http.NewRequest("GET", url, nil)

	reqest.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	reqest.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("Cache-Control", "max-age=0")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Referer", url)

	resp, err := client.Do(reqest)
	if err != nil {
		fmt.Println(url, err)
		return nil
	}

	defer resp.Body.Close()
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println(url, err)
			return nil
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	if reader != nil {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println(url, err)
			return nil
		}
		return body
	}
	return nil
}

func getHttpConsoleURL(doc *goquery.Document, nodeInfo string) string {
	var ret string
	doc.Find(nodeInfo).Each(func(i int, s *goquery.Selection) {
		tdAttr, exists := s.Attr("href")
		if exists &&strings.Contains(tdAttr, "get_proto_statistics") {
			ret = tdAttr
		}
	})
	return ret
}

func getHttpConsoleURLTest(doc *goquery.Document, nodeInfo string) [2]string {
	var ret [2]string
	doc.Find(nodeInfo).Each(func(i int, s *goquery.Selection) {
		if i != 0 || i != 1 {
			return
		}

		oumIndex := strconv.Itoa(i)
		var retTemp string
		subNodeInfo := nodeInfo + " .oum" + oumIndex + " td a"
		fmt.Printf("getHttpConsoleURLTest, subNodeinfo:%s\n", subNodeInfo)
		doc.Find(subNodeInfo).Each(func(i int, ss *goquery.Selection) {
			tdAttr, exists := ss.Attr("href")
			if exists &&strings.Contains(tdAttr, "get_proto_statistics") {
				retTemp = tdAttr
			}
		})
		ret[i] = retTemp
	})
	return ret
}

//GetAllMemdbSids 获取memdb所有子组的sid
func GetAllMemdbSids(inputURL *string) ([]string, error) {
	var ret []string
	doc, err := goquery.NewDocument(*inputURL)
	if err != nil {
		fmt.Print(err)
		return ret, err
	}

	doc.Find("body a").Each(func(i int, s *goquery.Selection) {
		attrName, exists := s.Attr("href")
		if exists && strings.Contains(attrName, "memdb.php?cid") {
			ret = append(ret, s.Text())
		}
	})
	return ret, err
}

// 执行http_console命令
func doExecute(httpConsoleUrl [2]string, command string) {
	//fmt.Printf(httpConsoleUrl + "\n") //查看获取到的http_console命令地址
	for i := 0; i < len(httpConsoleUrl); i++ {
		if httpConsoleUrl[i] == "" {
			continue
		}

		subipWithPort := strings.Split(httpConsoleUrl[i], "&")
		httpConsoleUrl[i] = subipWithPort[0] + command
		rsp := string(WebGetAction(httpConsoleUrl[i]))
		if strings.Contains(rsp, "+OK") {
			fmt.Printf(string(rsp) + "\n")
		}
	}
}

func DoExecuteHttpConsoleCMD(command string) {
	var memdbRootURL string = "http://www.xxx.com" //target url
	sids, errSids := GetAllMemdbSids(&memdbRootURL)
	if errSids != nil {
		fmt.Print(errSids)
		return
	}
	fmt.Println(sids)

	//遍历所有频道
	for _, v := range sids {
		OneMemdbURL := memdbRootURL + "?sid=" + v //拼凑单个语音组网页地址
		doc, err := goquery.NewDocument(OneMemdbURL)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("sid=%s:\n", v) //sid号
		doExecute(getHttpConsoleURLTest(doc, ".listtable"), command)
	}
}
