package main

import (
	"log"
	"net/http"
	"github.com/cheneylew/goutil/utils"
	"net/url"
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

func mainSignXiyu() {
	//loginHtml := getUrlRespHtml("http://oa.ehsy.com/login/Login.jsp?logintype=1")
	//utils.JJKPrintln(loginHtml)

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

	//签到
	//sign1 := url.Values{}
	//sign1.Set("signType","1")
	//sign1Html := utils.HTTPPostWithCookieCache("http://oa.ehsy.com/hrm/schedule/HrmScheduleSignXMLHTTP.jsp?t=0.31008359812174624",sign1)
	//utils.JJKPrintln(sign1Html)
	//
	////签退
	//sign2 := url.Values{}
	//sign2.Set("signType","2")
	//sign2Html := utils.HTTPPostWithCookieCache("http://oa.ehsy.com/hrm/schedule/HrmScheduleSignXMLHTTP.jsp?t=0.31008359812174624",sign2)
	//utils.JJKPrintln(sign2Html)

}