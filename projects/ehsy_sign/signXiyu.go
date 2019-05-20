package main

import (
	"log"
	"net/http"
	"github.com/cheneylew/goutil/utils"
	"net/url"
	"time"
	"path"
	"fmt"
	"strings"
)

func dbgPrintCurCookies() {
	var cookieNum int = len(utils.GCurCookies)
	log.Printf("cookieNum=%d", cookieNum)
	for i := 0; i < cookieNum; i++ {
		var curCk *http.Cookie = utils.GCurCookies[i]
		//log.Printf("curCk.Raw=%s", curCk.Raw)
		log.Printf("Cookie [%d]", i)
		log.Printf("Name\t=%s", curCk.Name)
		log.Printf("Value\t=%s", curCk.Value)
		log.Printf("Path\t=%s", curCk.Path)
		log.Printf("Domain\t=%s", curCk.Domain)
		log.Printf("Expires\t=%s", curCk.Expires)
		log.Printf("RawExpires=%s", curCk.RawExpires)
		log.Printf("MaxAge\t=%d", curCk.MaxAge)
		log.Printf("Secure\t=%t", curCk.Secure)
		log.Printf("HttpOnly=%t", curCk.HttpOnly)
		log.Printf("Raw\t=%s", curCk.Raw)
		log.Printf("Unparsed=%s", curCk.Unparsed)
	}
}

func testBaidu()  {
	//step1: access baidu url to get cookie BAIDUID
	log.Printf("======BAIDUID Cookie ======")
	var baiduMainUrl string = "http://www.baidu.com/"
	log.Printf("baiduMainUrl=%s", baiduMainUrl)
	respHtml := utils.HTTPGetWithCookieCache(baiduMainUrl)
	log.Printf("respHtml=%s", respHtml)
	dbgPrintCurCookies()

	//check cookie

	//step2: login, pass paras, extract resp cookie
	log.Printf("======login_token ======")
	//https://passport.baidu.com/v2/api/?getapi&class=login&tpl=mn&tangram=true
	var getapiUrl string = "https://passport.baidu.com/v2/api/?getapi&class=login&tpl=mn&tangram=true"
	var getApiRespHtml string = utils.HTTPGetWithCookieCache(getapiUrl)
	log.Printf("getApiRespHtml=%s", getApiRespHtml)
	dbgPrintCurCookies()
}

func login()  {
	//登陆
	vals := url.Values{}
	vals.Set("loginfile","/wui/theme/ecology8/page/login.jsp?templateId=4&logintype=1&gopage=")
	vals.Set("logintype","1")
	vals.Set("fontName","微软雅黑")
	vals.Set("message","")
	vals.Set("gopage","")
	vals.Set("formmethod","post")
	vals.Set("rnd","")
	vals.Set("serial","")
	vals.Set("username","")
	vals.Set("isie","false")
	vals.Set("islanguid","7")
	vals.Set("loginid","03527")
	vals.Set("userpassword","tough1988")
	vals.Set("submit","登录")
	utils.JJKPrintln(vals.Encode())
	postHtml := utils.HTTPPostWithCookieCache("http://oa.ehsy.com/login/VerifyLogin.jsp", vals)
	utils.JJKPrintln(postHtml)


	//首页
	getHtml := utils.HTTPGetWithCookieCache("http://oa.ehsy.com/wui/main.jsp?templateId=1")
	utils.JJKPrintln(getHtml)
}

func signIn()  {
	//签到
	sign1 := url.Values{}
	sign1.Set("signType","1")
	sign1Html := utils.HTTPPostWithCookieCache("http://oa.ehsy.com/hrm/schedule/HrmScheduleSignXMLHTTP.jsp?t=0.31008359812174624",sign1)
	utils.JJKPrintln(sign1Html)
	if strings.Contains(sign1Html, "您已成功签到") || strings.Contains(sign1Html, "签到（签退）时间") {
		utils.JJKPrintln("签到成功!")
		writeLog("签到成功!")
	} else {
		utils.JJKPrintln("签到失败!")
		writeLog("签到失败!")
	}
}

func signOut()  {
	//签退
	sign2 := url.Values{}
	sign2.Set("signType","2")
	sign2Html := utils.HTTPPostWithCookieCache("http://oa.ehsy.com/hrm/schedule/HrmScheduleSignXMLHTTP.jsp?t=0.31008359812174624",sign2)
	utils.JJKPrintln(sign2Html)
	if strings.Contains(sign2Html, "签到（签退）时间") || strings.Contains(sign2Html, "您已成功签退") {
		utils.JJKPrintln("签退成功!")
		writeLog("签退成功!")
	} else {
		utils.JJKPrintln("签退失败!")
		writeLog("签到失败!")
	}
}

func currentDate() string {
	return utils.JKDateNowStr()
}

func currentWeekDay() string {
	t := utils.JKStringToDate(currentDate())
	return t.Weekday().String()
}

func isWorkDay() bool {
	workdays := []string{
		"2018-12-29",
		"2019-02-02",
		"2019-02-03",
		"2019-09-29",
		"2019-10-12",
	}
	restDays := []string{
		"2018-12-31",
		"2019-01-01",
		"2019-02-04",
		"2019-02-05",
		"2019-02-06",
		"2019-02-07",
		"2019-02-08",
		"2019-02-09",
		"2019-02-10",
		"2019-04-05",
		"2019-05-01",
		"2019-06-07",
		"2019-09-13",
		"2019-10-01",
		"2019-10-02",
		"2019-10-03",
		"2019-10-04",
		"2019-10-05",
		"2019-10-06",
		"2019-10-07",
	}
	isExtWork := utils.InSlice(currentDate(), workdays)
	isExtRest := utils.InSlice(currentDate(), restDays)
	if isExtRest {
		return false
	}
	if isExtWork {
		return true
	}

	if currentWeekDay() == "Saturday" || currentWeekDay() == "Sunday" {
		return false
	} else {
		return true
	}
}

func writeLog(msg string)  {
	filePath := path.Join(utils.ExeDir(),"log.txt")
	if !utils.FileExist(filePath) {
		utils.FileWriteString(filePath,"")
	}
	fileContent := utils.FileReadAllString(filePath)
	fileContent += fmt.Sprintf("%s %s\r\n",utils.JKTimeNowStr(),  msg)
	utils.FileWriteString(filePath, fileContent)
}

func mainSignXiyu() {
	signalExit := make(chan int, 1)
	utils.CronJob("00 30 08 * * *", func() {
		if isWorkDay() {
			//休眠几分钟，每次不一样
			sleepCount := utils.RandomIntBetween(0,10)
			time.Sleep(time.Duration(sleepCount)*time.Minute)
			login()
			signIn()
		}
	})

	utils.CronJob("00 10 18 * * *", func() {
		if isWorkDay() {
			//休眠几分钟，每次不一样
			sleepCount := utils.RandomIntBetween(0,20)
			//周一 周二 周四 加班
			if 	currentWeekDay() == "Monday" ||
				currentWeekDay() == "Tuesday" ||
				currentWeekDay() == "Thursday" {
				sleepCount = 60*3+20+utils.RandomIntBetween(0,5)
			}
			time.Sleep(time.Duration(sleepCount)*time.Minute)
			login()
			signOut()
		}
	})

	utils.JJKPrintln("签到程序,已启动！")
	writeLog("签到程序,已启动！");
	<-signalExit
}