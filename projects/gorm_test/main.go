package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
	"time"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/go-redis/redis"
	"github.com/cheneylew/goutil/utils"
	"strings"
)

type Product struct {
	gorm.Model
	Code string
	Price uint
}

func goOrmFunc()  {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database :")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.First(&product, 1) // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)
}

func redisMain()  {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	client.Set("name", "cheneylew", time.Second*10)
	for {
		value := client.Get("name")
		utils.JJKPrintln(value.String())
		time.Sleep(time.Second * 1)
	}
}

var WORDS []string

func allWords() []string  {
	fpath := utils.ExeDirAppend("words.txt")
	timer := utils.JKTimer{}
	timer.Start()
	var str string
	if !utils.FileExist(fpath) {
		str = utils.HTTPGet("https://raw.githubusercontent.com/dwyl/english-words/master/words.txt")
		utils.FileWriteString(fpath, str)
	} else {
		str = utils.FileReadAllString(fpath)
	}
	timer.Record()
	result := strings.Split(str, "\n")
	timer.Record()
	return result
}

func radomWord() string {
	if len(WORDS) == 0 {
		WORDS = allWords()
	}
	index := utils.RandomIntBetween(0,len(WORDS))
	return WORDS[index]
}



func gormMysql()  {
	db, err := gorm.Open("mysql", "root:cnldj1988@tcp(localhost:13306)/cms?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		utils.JJKPrintln(err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&User{})

	timer := utils.JKTimer{}
	timer.Start()
	////插入1百万条数据
	//for i:=0;i<1000000 ; i++ {
	//	//插入
	//	username := radomWord()
	//	user := &User{
	//		Name:username,
	//		Email:username+"@163.com",
	//		Password:utils.RandomString(20),
	//		Birthday:"19880717",
	//	}
	//	db.Create(user)
	//}
	//timer.Record()

	//var users []*User
	//db.Where("email LIKE ?", "viz%").Find(&users)
	//timer.Record()
	//db.Where("name LIKE ?", "viz%").Find(&users)
	//timer.Record()

	v1, v2, err := utils.NewCaptchaBase64Str(100,40)
	utils.JJKPrintln(v1, v2, err)
	utils.FileWriteString("/Users/dejunliu/Desktop/a.txt", v2)
}

func main()  {
	gormMysql()

	utils.JJKPrintln("this is end!")
	select {
	}
}


