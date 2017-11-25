package main

import (
	"github.com/cheneylew/goutil/utils"
)

func main() {
	utils.JJKPrintln("hello world")

	p, n, v, e := utils.NewCaptcha("/Users/apple/Desktop",4,100,30)
	utils.JJKPrintln(p, n, v, e)

}