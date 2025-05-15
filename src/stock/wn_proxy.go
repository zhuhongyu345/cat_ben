package stock

import (
	"cat_ben/src/pkg/bizcall"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func getWBBasic(name string, mic string) (*WBBasicJson, error) {
	url := fmt.Sprintf("https://quotes-gw.webullfintech.com/api/quotes/ticker/getRealTimeBySymbol?exchangeCode=%s&symbol=%s", getWBCodeByMic(mic), name)
	headers := map[string]string{
		"accept":             "*/*",
		"accept-language":    "zh-CN,zh;q=0.9",
		"app":                "hk",
		"app-group":          "broker",
		"appid":              "wb_web_hk",
		"cache-control":      "no-cache",
		"device-type":        "Web",
		"did":                "fonyomm917n2dvrocefdw2e8vf5a9orf",
		"hl":                 "zh-hant",
		"origin":             "https://www.webull.hk",
		"os":                 "web",
		"osv":                "i9zh",
		"platform":           "web",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://www.webull.hk/quote/NYSE-AAMI",
		"reqid":              "tuxkq23vthr1vtqqvzcytrpwoj9kymq9",
		"sec-ch-ua":          "\"Chromium\";v=\"136\", \"Google Chrome\";v=\"136\", \"Not.A/Brand\";v=\"99\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "cross-site",
		"t_time":             "1747285183040",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
		"ver":                "1.0.0",
		"x-s":                "8d53c717ee702fc9463cae32eb2753e19ca9fe4fded37aa31b15f16a64d0c16c",
		"x-sv":               "xodp2vg9",
	}

	body, err := bizcall.GetJSONWithHeader(context.Background(), url, headers)
	if err != nil {
		log.Printf("wb basic call err:%s", err)
		return nil, err
	}
	var resp WBBasicJson
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("wb basic format err:%s", err)
		return nil, err
	}
	return &resp, nil
}

type WBBasicJson struct {
	TickerID             int    `json:"tickerId"`
	ExchangeID           int    `json:"exchangeId"`
	Type                 int    `json:"type"`
	SecType              []int  `json:"secType"`
	RegionID             int    `json:"regionId"`
	RegionCode           string `json:"regionCode"`
	CurrencyID           int    `json:"currencyId"`
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	DisSymbol            string `json:"disSymbol"`
	DisExchangeCode      string `json:"disExchangeCode"`
	ExchangeCode         string `json:"exchangeCode"`
	ListStatus           int    `json:"listStatus"`
	Template             string `json:"template"`
	DerivativeSupport    int    `json:"derivativeSupport"`
	TradeTime            string `json:"tradeTime"`
	Status               string `json:"status"`
	Close                string `json:"close"`
	Change               string `json:"change"`
	ChangeRatio          string `json:"changeRatio"`
	PPrice               string `json:"pPrice"`
	PChange              string `json:"pChange"`
	PChRatio             string `json:"pChRatio"`
	MarketValue          string `json:"marketValue"`
	Volume               string `json:"volume"`
	TurnoverRate         string `json:"turnoverRate"`
	TimeZone             string `json:"timeZone"`
	TzName               string `json:"tzName"`
	PreClose             string `json:"preClose"`
	Open                 string `json:"open"`
	High                 string `json:"high"`
	Low                  string `json:"low"`
	VibrateRatio         string `json:"vibrateRatio"`
	AvgVol10D            string `json:"avgVol10D"`
	AvgVol3M             string `json:"avgVol3M"`
	NegMarketValue       string `json:"negMarketValue"`
	Pe                   string `json:"pe"`
	ForwardPe            string `json:"forwardPe"`
	IndicatedPe          string `json:"indicatedPe"`
	PeTtm                string `json:"peTtm"`
	Eps                  string `json:"eps"`
	EpsTtm               string `json:"epsTtm"`
	Pb                   string `json:"pb"`
	TotalShares          string `json:"totalShares"`
	OutstandingShares    string `json:"outstandingShares"`
	FiftyTwoWkHigh       string `json:"fiftyTwoWkHigh"`
	FiftyTwoWkLow        string `json:"fiftyTwoWkLow"`
	Dividend             string `json:"dividend"`
	Yield                string `json:"yield"`
	BaSize               int    `json:"baSize"`
	NtvSize              int    `json:"ntvSize"`
	CurrencyCode         string `json:"currencyCode"`
	LotSize              string `json:"lotSize"`
	LatestDividendDate   string `json:"latestDividendDate"`
	LatestEarningsDate   string `json:"latestEarningsDate"`
	Ps                   string `json:"ps"`
	Bps                  string `json:"bps"`
	EstimateEarningsDate string `json:"estimateEarningsDate"`
	TradeStatus          string `json:"tradeStatus"`
}

func getWBCodeByMic(mic string) string {
	switch mic {
	case "XNYS":
		return "NYSE"
	case "XNGS", "XNMS", "XNCM":
		return "NASDAQ"
	case "XASE":
		return "AMEX"
	case "ARCX":
		return "NYSEARCA"
	case "BATS":
		return "BATS"
	case "OTCPK":
		return "OTCPK"
	default:
		log.Printf("unknown mic in wb:%s", mic)
		return mic
	}
}
