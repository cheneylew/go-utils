package utils

import (
	"github.com/astaxie/beego/httplib"
	"fmt"
	"net/http"
	"io/ioutil"
	"path"
	"strings"
	"os"
	"net/http/cookiejar"
	"log"
	"net/url"
)

var GCurCookies []*http.Cookie
var GCurCookieJar *cookiejar.Jar

//do init before all others
func init() {
	GCurCookies = nil
	//var err error;
	GCurCookieJar, _ = cookiejar.New(nil)
}


/*Get 访问一个URL*/
func HTTPGet(url string) string  {
	req := httplib.Get(url)
	str, err := req.String()
	if err != nil {
		Log(url, err)
	}
	return str
}
/*Post 方式访问一个URL*/
func HTTPPost(url string, params map[string]string) string  {
	req := httplib.Post(url)
	for k, v := range params {
		req.Param(k,v)
	}
	str, err := req.String()
	if err != nil {
		Log(url, err)
	}
	return str
}

//get url response html
func HTTPGetWithCookieCache(url string) string {
	log.Printf("getUrlRespHtml, url=%s", url)

	var respHtml string = ""

	httpClient := &http.Client{
		CheckRedirect: nil,
		Jar:           GCurCookieJar,
	}

	// httpResp, err := httpClient.Get("http://example.com")

	httpReq, err := http.NewRequest("GET", url, nil)
	//httpReq.Header.Add("If-None-Match", `W/"wyzzy"`)
	httpResp, err := httpClient.Do(httpReq)

	//httpResp, err := http.Get(url)
	//log.Printf("http.Get done")
	if err != nil {
		log.Printf("http get url=%s response error=%s\n", url, err.Error())
	}
	log.Printf("httpResp.Header=%s", httpResp.Header)
	log.Printf("httpResp.Status=%s", httpResp.Status)

	defer httpResp.Body.Close()

	body, errReadAll := ioutil.ReadAll(httpResp.Body)
	if errReadAll != nil {
		log.Printf("get response for url=%s got error=%s\n", url, errReadAll.Error())
	}

	GCurCookies = GCurCookieJar.Cookies(httpReq.URL)

	respHtml = string(body)
	return respHtml
}

func HTTPPostWithCookieCache(url string, values url.Values) string {
	log.Printf("getUrlRespHtml, url=%s", url)

	var respHtml string = ""

	httpClient := &http.Client{
		CheckRedirect: nil,
		Jar:           GCurCookieJar,
	}

	// httpResp, err := httpClient.Get("http://example.com")

	httpReq, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//httpReq.Header.Add("If-None-Match", `W/"wyzzy"`)
	httpResp, err := httpClient.Do(httpReq)

	//httpResp, err := http.Get(url)
	//log.Printf("http.Get done")
	if err != nil {
		log.Printf("http get url=%s response error=%s\n", url, err.Error())
	}
	log.Printf("httpResp.Header=%s", httpResp.Header)
	log.Printf("httpResp.Status=%s", httpResp.Status)

	defer httpResp.Body.Close()

	body, errReadAll := ioutil.ReadAll(httpResp.Body)
	if errReadAll != nil {
		log.Printf("get response for url=%s got error=%s\n", url, errReadAll.Error())
	}

	GCurCookies = GCurCookieJar.Cookies(httpReq.URL)

	respHtml = string(body)
	return respHtml
}


func DJMapToHttpGetParams(filter map[string]string) string {
	result := ""
	for key, value := range filter {
		if len(value) > 0 {
			result +=fmt.Sprintf("&%s=%s",Trim(key),Trim(value))
		}
	}

	return result
}

func DJDownloadImageToDefaultDir(url string) string {
	imageDirPath := MakeDir(fmt.Sprintf("downloads/images/%s/",JKDateNowStr()))
	imagePath := DJDownloadImage(url,imageDirPath)
	return imagePath
}

func DJDownloadImage(url string, relativeDirPath string) string {
	response,err := http.Get(url)
	if err != nil {
		fmt.Println("download image error:",err.Error())
	}
	defer response.Body.Close()
	tp := response.Header.Get("Content-Type")
	tpArray := strings.Split(tp,"/")
	fileType := ""
	if len(tpArray) > 0{
		fileType = "."+tpArray[len(tpArray)-1]
	}
	imageBytes,imagErr := ioutil.ReadAll(response.Body)
	if imagErr != nil {
		fmt.Println("read image bytes error:",imagErr.Error())
	}

	imagePath := path.Join(relativeDirPath,JKIntToStr(int(JKTimeNowStamp()))+"-"+JKIntToStr(RandomInt(10000))+fileType)
	f,e := os.Create(imagePath)
	if e != nil {
		fmt.Println("create image file error :",e.Error())
	}
	defer f.Close()
	_,e1 := f.Write(imageBytes)
	if e1 != nil {
		fmt.Println("write image bytes error:",e1.Error())
	}
	return imagePath
}