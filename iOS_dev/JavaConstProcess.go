package main

import (
	"github.com/cheneylew/goutil/utils"
	"strings"
	"path"
	"fmt"
)

type Const struct {
	Key string
	Value string
	Comment string
	ClassName string
}

func filePath(name string) string {
	dir := "/Users/dejunliu/Desktop/"
	return path.Join(dir, name)
}

func strToConst(str string) *Const {
	if strings.Contains(str, "=") {
		a := strings.Split(str,"=")
		if len(a) == 2 {
			var key,value,comment string
			key = utils.Trim(strings.Replace(a[0],"String","",-1))
			valTMP := utils.Trim(strings.Replace(a[1],"","",-1))
			com := strings.Split(valTMP,";")
			value = utils.Trim(strings.Replace(com[0],"\"","",-1))
			if len(com) == 2 {
				comment = utils.Trim(strings.Replace(com[1],"//","",-1))
			}
			return &Const{
				Key:key,
				Value:value,
				Comment:comment,
			}
		}
	}


	return nil
}

func JavaConstProcess()  {

	javaConst := utils.FileReadAllString(filePath("BizConstant.java"))
	rows := strings.Split(javaConst,"\n")

	results := make([]*Const,0)
	javaInterface := ""
	for _, row := range rows {
		if strings.Contains(row, "interface") {
			t := strings.Split(row,"interface")
			t1 := strings.Split(t[1],"{")
			t2 := utils.Trim(t1[0])
			javaInterface = t2
			utils.JJKPrintln(utils.SnakeString(javaInterface), t2)
		}
		jConst := strToConst(row)
		if jConst != nil {
			jConst.ClassName = javaInterface
			results = append(results, jConst)
		}

	}

	strs := ""
	lastClassName := ""
	for _, result := range results {
		enter := ""
		if len(lastClassName) == 0 {
			enter = ""
		}
		if lastClassName != result.ClassName {
			enter = "\n"
		}

		lastClassName = result.ClassName
		strs += fmt.Sprintf("%s#define K%s_%s\t\t@\"%s\"\t//%s\n",
			enter,
			result.ClassName,
			result.Key,
			result.Value,
			result.Comment)
	}
	utils.FileWriteString(filePath("a.pch"), strs)


}


