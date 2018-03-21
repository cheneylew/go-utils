package main

import (
	"github.com/cheneylew/goutil/utils"
	"time"
)

func main() {

	for i:=0; i<10 ; i++ {
		utils.JJKPrintln(i)
	}

	singal := make(chan int)
	utils.JJKPrintln("Started!")

	go func(j chan int) {
		i := 10
		for i > 0 {
			i --
			time.Sleep(time.Second*2)
			utils.JJKPrintln("ok!")
		}
		j <- 1
	}(singal)

	<- singal
	utils.JJKPrintln("End!")
}