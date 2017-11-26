package main

import (
	"bufio"
	"strings"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"github.com/cheneylew/goutil/utils"
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

func main() {
	ExampleScrape()
}