package main

import "github.com/cheneylew/goutil/utils"

func main() {
	timer := utils.JKTimer{}
	timer.Start()
	count := 10000*1000
	slice := make([]string, count)
	for i:=0;i<count ; i++ {
		slice = append(slice, utils.RandomString(32))
	}
	timer.Record()
	for _, value := range slice {
		if len(value) == 32 {
			
		}
	}
	timer.Record()

	utils.JJKPrintln("End!")
	select {
	}
}