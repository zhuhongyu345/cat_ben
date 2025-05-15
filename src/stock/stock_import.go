package stock

import (
	"cat_ben/src/db"
	"cat_ben/src/pkg/bizcall"
	"context"
	"encoding/json"
	"log"
	"strings"
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
	temp := GetFromNyse("")
	unknown := make([]string, 0)
	stos := make([]*db.Sto, 0)
	for _, t := range temp {
		if strings.Contains(t.NormalizedTicker, ".") || strings.Contains(t.NormalizedTicker, ":") {
			continue
		}
		sType := 0
		if t.InstrumentType == "EXCHANGE_TRADED_FUND" || t.InstrumentType == "INDEX" {
			sType = 2
		} else if t.InstrumentType == "COMMON_STOCK" || t.InstrumentType == "DEPOSITORY_RECEIPT" || t.InstrumentType == "PREFERRED_STOCK" {
			sType = 1
		} else if t.InstrumentType == "RIGHT" || t.InstrumentType == "UNIT" || t.InstrumentType == "REIT" ||
			t.InstrumentType == "LIMITED_PARTNERSHIP" || t.InstrumentType == "WARRANT" || t.InstrumentType == "CLOSED_END_FUND" ||
			t.InstrumentType == "EXCHANGE_TRADED_NOTE" || t.InstrumentType == "NOTE" || t.InstrumentType == "TEST" ||
			t.InstrumentType == "TRUST" || t.InstrumentType == "UNITS_OF_BENEFICIAL_INTEREST" || t.InstrumentType == "BOND" {
			continue
		} else {
			log.Printf("%+v", t)
			unknown = append(unknown, t.InstrumentType)
		}
		stos = append(stos, &db.Sto{
			Name: t.NormalizedTicker,
			Type: sType,
			Mic:  t.MicCode,
		})
	}
	log.Printf("%+v", unknown)
	if len(unknown) > 0 {
		panic(unknown)
	}
	_ = db.DeleteALL()
	for _, sto := range stos {
		err2 := db.CreateStos([]*db.Sto{sto})
		if err2 != nil {
			log.Printf("%s", err2)
			panic(err2)
		}
	}

	time.Sleep(time.Second * 10)
}

func GetFromNyse(name string) []*NyseUnit {
	serverURL := "https://www.nyse.com/api/quotes/filter"
	pageSize := 50000
	param := map[string]interface{}{
		"pageNumber":        1,
		"sortColumn":        "NORMALIZED_TICKER",
		"sortOrder":         "ASC",
		"maxResultsPerPage": pageSize,
	}
	if name != "" {
		param["filterToken"] = name
	}
	resp, _ := bizcall.PostJSONWithHeader(context.TODO(), serverURL, param, nil)
	temp := make([]*NyseUnit, 0)
	_ = json.Unmarshal(resp, &temp)
	return temp
}
