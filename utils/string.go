package utils

import (
	"strconv"
	"regexp"
	"strings"
	"encoding/base64"
	"html/template"
	"bytes"
	"github.com/astaxie/beego"
)


func JKStrToInt(str string) int {
	it,_ := strconv.Atoi(str)
	return it
}

func JKStrToUInt8(str string) uint8 {
	it,_ := strconv.Atoi(str)
	return uint8(it)
}

func JKStrToInt64(str string) int64 {
	it,_ := strconv.Atoi(str)
	return int64(it)
}

func JKIntToStr(i int) string {
	return strconv.Itoa(i)
}

func JKHTMLEscape(str string) string  {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src := re.ReplaceAllStringFunc(str, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

func UpperFirstChar(str string) string {
	if len(str) > 0 {
		f := str[:1]
		return strings.ToUpper(f) + str[1:]
	}

	return ""
}

func LowerFirstChar(str string) string {
	if len(str) > 0 {
		f := str[:1]
		return strings.ToLower(f) + str[1:]
	}

	return ""
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) string {
	b, _ := base64.StdEncoding.DecodeString(str)
	return string(b)
}

func TemplateParams() map[string]interface{} {
	return make(map[string]interface{}, 0)
}

func Template(templateStr string, params map[string]interface{}) string {

	t := template.Must(template.New("tpl").Funcs(template.FuncMap{
		"Equal":Equal,
		"InSlice":InSlice,
		"ToStr":ToString,
		"MapGet":beego.MapGet,
	}).Parse(templateStr))

	buf := bytes.NewBufferString("")
	t.Execute(buf, params)

	return buf.String()
}