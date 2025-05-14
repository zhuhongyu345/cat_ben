package server

import (
	"cat_ben/src/chromedriver"
	"cat_ben/src/stock"
	"time"
)

var doing = false

func FlushTask() {
	for {
		time.Sleep(time.Hour * 1)
		doing = true
		chromedriver.GetTokenAndSave()
		stock.FlushBasic("1", "")
		doing = false
	}
}
