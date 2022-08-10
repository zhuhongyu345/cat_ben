package stock

import (
	"cat_ben/src/db"
	"log"
)

func Search(name, hlLow, hlHight, peHigh, peLow, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType string) []*db.Sto {
	stocks, err := db.Search(name, hlLow, hlHight, peHigh, peLow, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType)
	if err != nil {
		log.Printf("db.Search err:%s", err)
	}
	return stocks
}
