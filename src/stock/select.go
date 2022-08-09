package stock

import (
	"cat_ben/src/db"
	"log"
)

func Search(hlLow, hlHight, pe, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType string) []*db.Sto {
	stocks, err := db.Search(hlLow, hlHight, pe, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType)
	if err != nil {
		log.Printf("db.Search err:%s", err)
	}
	return stocks
}
