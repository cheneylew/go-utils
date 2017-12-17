package stock

import (
	"log"
	"github.com/cheneylew/goutil/utils"
)

func CronMain()  {
	utils.CronJob("1-10/2 * * * * *", func() {
		log.Println("ok")
	})
}
