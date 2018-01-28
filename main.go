package main

import (
	"bufio"
	"strings"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"github.com/cheneylew/goutil/utils"
	"github.com/astaxie/beego/cache"
	"time"
)

func ExampleScrape() {
	doc, err := goquery.NewDocument("https://www.haiyinhui.com")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".fl").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		//band := s.Find("a").Text()
		//title := s.Find("i").Text()
		//fmt.Printf("Review %d: %s - %s\n", i, band, title)
		utils.JJKPrintln(s.Text())
	})
}

func Reader()  {
	reader := bufio.NewReader(strings.NewReader("http://studygolang.com. \nIt is the home of gophers"))
	line, _ := reader.ReadBytes('\n')
	fmt.Printf("the line:%s\n", line)
	// 这里可以换上任意的 bufio 的 Read/Write 操作
	n, _ := reader.ReadBytes('\n')
	fmt.Printf("the line:%s\n", line)
	fmt.Println(string(n))
}

type User struct {
	Name string
	Age int
}

func main() {
	//ExampleScrape()
	//utils.QRCodePNGWriteFile("http://www.baidu.com/","/Users/apple/Desktop/a.png")
	//s, e := utils.QRCodeDecode("/Users/apple/Desktop/a.png")
	//utils.JJKPrintln(s,e)

	//内存缓存
	if false {
		bm, err := cache.NewCache("memory", `{"interval":60}`)
		if err != nil {
			utils.JJKPrintln(err)
		} else {
			bm.Put("name","cheneylew",time.Second * 10)
			ticker := time.NewTicker(time.Second * 1)
			go func() {
				for _ = range ticker.C {
					utils.JJKPrintln(bm.Get("name"))
				}
			}()
		}
	}

	//并发任务队列
	if false {
		var params []interface{}
		for i:=0; i< 5; i++ {
			params = append(params, fmt.Sprintf("params %d", i))
		}

		utils.QueueTask(2, params, func(idx int, param interface{}) {
			utils.JJKPrintln(param)
			time.Sleep(time.Second * 2)
		})

	}

	json := utils.HTTPGet("https://bx.in.th/api/orderbook/?pairing=21")
	utils.JJKPrintln(json)

	utils.JJKPrintln("all task finished!")
	select {}
}