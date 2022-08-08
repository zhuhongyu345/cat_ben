package option

import (
	"cat_ben/src/pkg/bizcall"
	"context"
	"encoding/json"
	"log"
	"sort"
	"strconv"
)

func GetOptionChain(id int64) map[string]interface{} {
	serverURL := "https://quotes-gw.webullfintech.com/api/quote/option/strategy/list"
	header := map[string]string{
		"access_token": "dc_cn1.180947de508-a5ae8e46b41b4332b6e3b5b7bf4b46d5",
		"content-type": "application/json",
		"device-type":  "Web",
		"did":          "9f315827fa764ff6bd5bf619d3f0fb5f",
		"hl":           "zh",
		"locale":       "zh_CN",
		"os":           "web",
		"osv":          "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"ver":          "3.39.19",
	}
	param := map[string]interface{}{
		"tickerId":        id,
		"count":           50,
		"direction":       "all",
		"expireCycle":     []int{3, 2, 4},
		"type":            0,
		"quoteMultiplier": 100,
	}
	resp, err := bizcall.PostJSONWithHeader(context.TODO(), serverURL, param, header)
	if err != nil {
		log.Printf("%s", err)
	}
	//log.Printf("chain:%s", string(resp))
	var chain ChainWebull
	json.Unmarshal(resp, &chain)

	if len(chain.ExpireDateList) == 0 {
		log.Printf("err")
	}
	list := chain.ExpireDateList[0]
	//当前标的价格。
	nowPrice, _ := strconv.ParseFloat(chain.Close, 64)
	//if len(chain.AskList) > 0 && len(chain.BidList) > 0 {
	//	askP, _ := strconv.ParseFloat(chain.AskList[0].Price, 64)
	//	bidP, _ := strconv.ParseFloat(chain.BidList[0].Price, 64)
	//	nowPrice = (askP + bidP) / 2
	//}
	//行权价：横坐标
	strikes := make([]float64, 0)
	for _, data := range list.Data {
		float, _ := strconv.ParseFloat(data.StrikePrice, 64)
		if data.Direction == "call" {
			strikes = append(strikes, float)
		}
	}
	sort.Float64s(strikes)
	//resp
	respData := map[string]interface{}{
		"nowPrice": nowPrice,
		"strikes":  strikes,
	}

	for idx, list := range chain.ExpireDateList {
		if idx == 0 {
			respData["char1"] = GetCharFromList(nowPrice, strikes, list)
		}
		if idx == 1 {
			respData["char2"] = GetCharFromList(nowPrice, strikes, list)
		}
		if idx > 0 && list.From.Weekly == 0 {
			respData["charMonth"] = GetCharFromList(nowPrice, strikes, list)
			break
		}
	}
	return respData
}

func GetCharFromList(nowPrice float64, strikes []float64, list ExpireDateList) map[string]interface{} {
	openInterestCall := make([]int, 0)
	openInterestCallMoney := make([]float64, 0)
	openInterestPut := make([]int, 0)
	openInterestPutMoney := make([]float64, 0)
	deltas := make([]float64, 0)
	chain := list.Data
	for _, s := range strikes {
		rate := float64(0)
		for _, data := range chain {
			gamma, _ := strconv.ParseFloat(data.Gamma, 64)
			delta, _ := strconv.ParseFloat(data.Delta, 64)
			sPrice, _ := strconv.ParseFloat(data.StrikePrice, 64)
			diffPrice := s - nowPrice
			money := (0.5*gamma*diffPrice*diffPrice + delta*diffPrice) * float64(data.OpenInterest)
			rate += money
			if sPrice == s {
				nowPrice := float64(0)
				if len(data.AskList) > 0 && len(data.BidList) > 0 {
					askP, _ := strconv.ParseFloat(data.AskList[0].Price, 64)
					bidP, _ := strconv.ParseFloat(data.BidList[0].Price, 64)
					nowPrice = (askP + bidP) / 2
				}
				if data.Direction == "call" {
					openInterestCall = append(openInterestCall, data.OpenInterest)
					openInterestCallMoney = append(openInterestCallMoney, nowPrice*float64(data.OpenInterest))
				} else if data.Direction == "put" {
					openInterestPut = append(openInterestPut, -data.OpenInterest)
					openInterestPutMoney = append(openInterestPutMoney, -nowPrice*float64(data.OpenInterest))

				}
			}
		}
		deltas = append(deltas, rate/10)
	}
	resp := map[string]interface{}{
		"openInterestCall":      openInterestCall,
		"openInterestPut":       openInterestPut,
		"openInterestCallMoney": openInterestCallMoney,
		"openInterestPutMoney":  openInterestPutMoney,
		"deltas":                deltas,
	}
	return resp
}

type KWebull struct {
	TickerID       int      `json:"tickerId"`
	TimeZone       string   `json:"timeZone"`
	PreClose       string   `json:"preClose"`
	HasMore        int      `json:"hasMore"`
	ExchangeStatus bool     `json:"exchangeStatus"`
	Data           []string `json:"data"`
	CleanTime      int64    `json:"cleanTime"`
	CleanDuration  int      `json:"cleanDuration"`
}

type ChainWebull struct {
	TickerID        int              `json:"tickerId"`
	Name            string           `json:"name"`
	DisSymbol       string           `json:"disSymbol"`
	ExchangeCode    string           `json:"exchangeCode"`
	DisExchangeCode string           `json:"disExchangeCode"`
	Close           string           `json:"close"`
	PreClose        string           `json:"preClose"`
	Volume          string           `json:"volume"`
	Open            string           `json:"open"`
	High            string           `json:"high"`
	Low             string           `json:"low"`
	TickerType      int              `json:"tickerType"`
	AskList         []AskList        `json:"askList"`
	BidList         []BidList        `json:"bidList"`
	Change          string           `json:"change"`
	ChangeRatio     string           `json:"changeRatio"`
	NtvSize         int              `json:"ntvSize"`
	Template        string           `json:"template"`
	ExpireDateList  []ExpireDateList `json:"expireDateList"`
	Vol1Y           string           `json:"vol1y"`
}

type From struct {
	Date            string `json:"date"`
	Days            int    `json:"days"`
	Weekly          int    `json:"weekly"`
	QuoteMultiplier int    `json:"quoteMultiplier"`
	QuoteLotSize    int    `json:"quoteLotSize"`
	UnSymbol        string `json:"unSymbol"`
}
type AskList struct {
	Price   string `json:"price"`
	Volume  string `json:"volume"`
	QuoteEx string `json:"quoteEx"`
}
type BidList struct {
	Price   string `json:"price"`
	Volume  string `json:"volume"`
	QuoteEx string `json:"quoteEx"`
}
type Data struct {
	Open             string    `json:"open"`
	High             string    `json:"high"`
	Low              string    `json:"low"`
	StrikePrice      string    `json:"strikePrice"`
	PreClose         string    `json:"preClose"`
	OpenInterest     int       `json:"openInterest"`
	Volume           string    `json:"volume"`
	LatestPriceVol   string    `json:"latestPriceVol"`
	Delta            string    `json:"delta"`
	Vega             string    `json:"vega"`
	ImpVol           string    `json:"impVol"`
	Gamma            string    `json:"gamma"`
	Theta            string    `json:"theta"`
	Rho              string    `json:"rho"`
	Close            string    `json:"close"`
	Change           string    `json:"change"`
	ChangeRatio      string    `json:"changeRatio"`
	ExpireDate       string    `json:"expireDate"`
	TickerID         int       `json:"tickerId"`
	BelongTickerID   int       `json:"belongTickerId"`
	OpenIntChange    int       `json:"openIntChange"`
	ActiveLevel      int       `json:"activeLevel"`
	Cycle            int       `json:"cycle"`
	Weekly           int       `json:"weekly"`
	ExecutionType    string    `json:"executionType"`
	Direction        string    `json:"direction"`
	DerivativeStatus int       `json:"derivativeStatus"`
	CurrencyID       int       `json:"currencyId"`
	RegionID         int       `json:"regionId"`
	ExchangeID       int       `json:"exchangeId"`
	Symbol           string    `json:"symbol"`
	UnSymbol         string    `json:"unSymbol"`
	AskList          []AskList `json:"askList"`
	BidList          []BidList `json:"bidList"`
	QuoteMultiplier  int       `json:"quoteMultiplier"`
	QuoteLotSize     int       `json:"quoteLotSize"`
}
type Call struct {
	Option  int    `json:"option"`
	Side    string `json:"side"`
	Gravity int    `json:"gravity"`
}
type Put struct {
	Option  int    `json:"option"`
	Side    string `json:"side"`
	Gravity int    `json:"gravity"`
}

type ExpireDateList struct {
	From From   `json:"from"`
	Data []Data `json:"data"`
}
