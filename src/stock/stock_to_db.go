package stock

import (
	"cat_ben/src/db"
	"cat_ben/src/pkg/bizcall"
	"context"
	"encoding/json"
	"log"
	"time"
)

type NyseUnit struct {
	ExchangeID           string `json:"exchangeId"`
	InstrumentName       string `json:"instrumentName"`
	InstrumentType       string `json:"instrumentType"`
	MicCode              string `json:"micCode"`
	NormalizedTicker     string `json:"normalizedTicker"`
	SymbolEsignalTicker  string `json:"symbolEsignalTicker"`
	SymbolExchangeTicker string `json:"symbolExchangeTicker"`
	SymbolTicker         string `json:"symbolTicker"`
	Total                int    `json:"total"`
	URL                  string `json:"url"`
}

func allStockToDB() {
	serverURL := "https://www.nyse.com/api/quotes/filter"
	pageSize := 100
	param := map[string]interface{}{
		"pageNumber":        1,
		"sortColumn":        "NORMALIZED_TICKER",
		"sortOrder":         "ASC",
		"maxResultsPerPage": pageSize,
	}
	header, err := bizcall.PostJSONWithHeader(context.TODO(), serverURL, param, nil)
	if err != nil {
		log.Printf("%s", err)
	}
	log.Printf("%s", header)
	page := make([]*NyseUnit, 0)
	_ = json.Unmarshal(header, &page)
	total := page[0].Total
	for i := 1; i < total/pageSize; i++ {
		param["pageNumber"] = i
		resp, _ := bizcall.PostJSONWithHeader(context.TODO(), serverURL, param, nil)
		temp := make([]*NyseUnit, 0)
		_ = json.Unmarshal(resp, &temp)
		stos := make([]*db.Sto, 0)
		for _, t := range temp {
			sType := 0
			if t.InstrumentType == "EXCHANGE_TRADED_FUND" {
				sType = 2
			} else if t.InstrumentType == "COMMON_STOCK" {
				sType = 1
			} else {
				continue
			}
			stos = append(stos, &db.Sto{
				Name: t.NormalizedTicker,
				Type: sType,
			})
		}
		err2 := db.CreateStos(stos)
		if err2 != nil {
			log.Printf("%s", err2)
		}
	}
	time.Sleep(time.Second * 100)
}
