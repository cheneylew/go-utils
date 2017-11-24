package utils

import (
	"github.com/astaxie/beego/httplib"
	"fmt"
	"net/http"
	"io/ioutil"
	"path"
	"strings"
	"os"
)

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

	imagePath := path.Join(relativeDirPath,JKIntToStr(int(JKTimeNowStamp()))+"-"+JKIntToStr(JKRandInt(10000))+fileType)
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