package main

import (
	"github.com/cheneylew/goutil/utils"
	"net/http"
	"io"
	"log"
	"time"
	"encoding/json"
)


type User struct {
	Id string
	Balance float64
}

func timerStart() {
	t1 := time.NewTicker(time.Millisecond * 3000)
	go func() {
		for range t1.C {
			//params := make(map[string]string)
			//params["username"] = "cheneylew"
			//params["password"] = "111111"
			//response := utils.HTTPPost("http://localhost:12345/json", params)
			//utils.JJKPrintln(response)
		}
	}()
}

func HttpServerStart() {
	utils.JJKPrintln("http server started!")
	timerStart()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		username := request.PostFormValue("username")
		password := request.PostFormValue("password")
		utils.JJKPrintln(username, password)
		io.WriteString(writer, "hello, world!\n")
	})

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "hello, world!\n")
	})

	http.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
		u := User{Id: "US123", Balance: 8}
		json.NewEncoder(writer).Encode(u)
	})

	err := http.ListenAndServe(":12345", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	utils.JJKPrintln("http server end!")
}
