package option

import (
	"cat_ben/src/pkg/bizcall"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// DataFromUws /*
func DataFromUws() {
	serverURL := fmt.Sprint("https://phx.unusualwhales.com/api/option_trades_v2" +
		"?limit=10&newer_than=2022-05-31T14:45:59.943Z&min_premium=100000&tags[]=ask_side" +
		"&tags[]=aviation&tags[]=bearish&tags[]=bid_side&tags[]=biotech&tags[]=bullish" +
		"&tags[]=china&tags[]=dividend&tags[]=earnings_this_week&tags[]=earnings_next_week&tags[]=expires_within_week" +
		"&tags[]=index&tags[]=leap&tags[]=expires_more_than_month&tags[]=oil" +
		"&tags[]=semi&tags[]=heavily_shorted&tags[]=spac&tags[]=volatility&tags[]=weed&order=Time&ticker_symbol=TSLA&min_dte=0",
	)
	header := map[string]string{
		"authority":       "phx.unusualwhales.com",
		"accept":          "application/json, text/plain",
		"accept-language": "zh-CN,zh;q=0.9",
		"authorization":   "Bearer w2VYpnRKiQaRM6x5T8Xihzer1_HNOTo06o-I2b11bHtnaWj49TdFnD1HniPvLFbn",
		"origin":          "https://unusualwhales.com",
		"referer":         "https://unusualwhales.com/",
		"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	}

	resp, err := bizcall.GetJSONWithHeader(context.TODO(), serverURL, header)
	if err != nil {
		panic(err)
	}
	var result UwsResult
	json.Unmarshal(resp, &result)
	log.Printf("%+v", result)

}

type UwsResult struct {
	Data []DataUws `json:"data"`
}
type DataUws struct {
	AlertScore              string        `json:"alert_score"`
	Canceled                bool          `json:"canceled"`
	CreatedAt               time.Time     `json:"created_at"`
	Delta                   string        `json:"delta"`
	EwmaNbboAsk             string        `json:"ewma_nbbo_ask"`
	EwmaNbboBid             string        `json:"ewma_nbbo_bid"`
	ExecutedAt              time.Time     `json:"executed_at"`
	Gamma                   string        `json:"gamma"`
	ID                      string        `json:"id"`
	ImpliedVolatility       string        `json:"implied_volatility"`
	MarketCenterLocate      int           `json:"market_center_locate"`
	NbboAsk                 string        `json:"nbbo_ask"`
	NbboBid                 string        `json:"nbbo_bid"`
	OpenInterest            int           `json:"open_interest"`
	OptionChainID           string        `json:"option_chain_id"`
	Price                   string        `json:"price"`
	ReportFlags             []interface{} `json:"report_flags"`
	Rho                     string        `json:"rho"`
	Sector                  string        `json:"sector"`
	Size                    int           `json:"size"`
	Strike                  string        `json:"strike"`
	Tags                    []string      `json:"tags"`
	Theo                    string        `json:"theo"`
	Theta                   string        `json:"theta"`
	UnderlyingPrice         string        `json:"underlying_price"`
	UnderlyingSymbol        string        `json:"underlying_symbol"`
	UpstreamConditionDetail string        `json:"upstream_condition_detail"`
	Vega                    string        `json:"vega"`
	Volume                  int           `json:"volume"`
}
