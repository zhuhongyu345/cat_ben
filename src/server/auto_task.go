package server

import (
	"cat_ben/src/chromedriver"
	"cat_ben/src/stock"
	"time"
)

func FlushTask() {
	for {
		time.Sleep(time.Hour * 1)
		chromedriver.GetTokenAndSave()
		now := time.Now() // 获取当前时间（本地时区）

		// 获取小时（24小时制）和分钟
		hour := now.Hour()
		// 逻辑判断：小时 > 9 或者 小时 == 9 且分钟 >= 0
		if hour >= 9 && hour <= 18 {
			stock.FlushBasic("0", "2")
			stock.FlushBasic("1", "1")
		} else {
			stock.FlushBasic("1", "1")
		}
	}
}
