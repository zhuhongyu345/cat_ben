package stock

import (
	"cat_ben/src/db"
	"cat_ben/src/pkg/bizcall"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
)

// from https://xueqiu.com/S/INTC
func FlushBasic(hard string) {
	log.Println("flush start")
	stocks, _ := db.GetAllStockFromDB(hard)
	for _, stock := range stocks {
		pe, yield, chn, price, h52, l52, liangbi, shizhi, huanshoulv, err := getDetailFromXQ(stock.Name)
		if err != nil {
			continue
		}
		hl := (price - l52) / (h52 - l52)
		if h52 == l52 {
			hl = 1
		}
		hl = math.Round(hl*10000) / 10000
		err = db.UpdateByID(stock.ID, pe, yield, chn, price, h52, l52, hl, liangbi, shizhi, huanshoulv)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("flush finish")
}

func getDetailFromXQ(name string) (float64, float64, string, float64, float64, float64, float64, float64, float64, error) {
	token, _ := db.GetValue("xueqiu_token")
	header := map[string]string{
		"cookie":     token,
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	}
	resp, err := bizcall.GetJSONWithHeader(context.Background(), fmt.Sprintf("https://stock.xueqiu.com/v5/stock/quote.json?symbol=%s&extend=detail", name), header)
	if err != nil {
		log.Println(err)
		return 0, 0, "", 0, 0, 0, 0, 0, 0, err
	}
	var respXQ JSONDataXueQiuBasic
	json.Unmarshal(resp, &respXQ)
	log.Println(respXQ)
	ttm := respXQ.Data.Quote.PeTtm
	yield := respXQ.Data.Quote.DividendYield
	nameCHN := respXQ.Data.Quote.Name
	price := respXQ.Data.Quote.Current
	h52 := respXQ.Data.Quote.High52W
	l52 := respXQ.Data.Quote.Low52W
	liangbi := respXQ.Data.Quote.VolumeRatio
	shizhi := (respXQ.Data.Quote.MarketCapital) / float64(100000000)
	huanshou := respXQ.Data.Quote.TurnoverRate
	return ttm, yield, nameCHN, price, h52, l52, liangbi, shizhi, huanshou, nil
}

type JSONDataXueQiuBasic struct {
	Data Data `json:"data"`
}

type Quote struct {
	CurrentExt          float64     `json:"current_ext"`
	Symbol              string      `json:"symbol"`
	High52W             float64     `json:"high52w"`
	PercentExt          float64     `json:"percent_ext"`
	Delayed             int         `json:"delayed"`
	Type                int         `json:"type"`
	TickSize            float64     `json:"tick_size"`
	FloatShares         interface{} `json:"float_shares"`
	High                float64     `json:"high"`
	FloatMarketCapital  interface{} `json:"float_market_capital"`
	TimestampExt        int64       `json:"timestamp_ext"`
	LotSize             int         `json:"lot_size"`
	LockSet             int         `json:"lock_set"`
	Chg                 float64     `json:"chg"`
	Eps                 float64     `json:"eps"`
	LastClose           float64     `json:"last_close"`
	ProfitFour          int64       `json:"profit_four"`
	Volume              int         `json:"volume"`
	VolumeRatio         float64     `json:"volume_ratio"`
	ProfitForecast      int64       `json:"profit_forecast"`
	TurnoverRate        float64     `json:"turnover_rate"`
	Low52W              float64     `json:"low52w"`
	Name                string      `json:"name"`
	Exchange            string      `json:"exchange"`
	PeForecast          float64     `json:"pe_forecast"`
	TotalShares         int64       `json:"total_shares"`
	Status              int         `json:"status"`
	Code                string      `json:"code"`
	GoodwillInNetAssets float64     `json:"goodwill_in_net_assets"`
	AvgPrice            float64     `json:"avg_price"`
	Percent             float64     `json:"percent"`
	Psr                 float64     `json:"psr"`
	Amplitude           float64     `json:"amplitude"`
	Current             float64     `json:"current"`
	CurrentYearPercent  float64     `json:"current_year_percent"`
	IssueDate           int64       `json:"issue_date"`
	SubType             string      `json:"sub_type"`
	Low                 float64     `json:"low"`
	MarketCapital       float64     `json:"market_capital"`
	ShareholderFunds    int64       `json:"shareholder_funds"`
	Dividend            float64     `json:"dividend"`
	DividendYield       float64     `json:"dividend_yield"`
	Currency            string      `json:"currency"`
	ChgExt              float64     `json:"chg_ext"`
	Navps               float64     `json:"navps"`
	Profit              int64       `json:"profit"`
	Beta                interface{} `json:"beta"`
	Timestamp           int64       `json:"timestamp"`
	PeLyr               float64     `json:"pe_lyr"`
	Amount              int         `json:"amount"`
	PledgeRatio         interface{} `json:"pledge_ratio"`
	ShortRatio          interface{} `json:"short_ratio"`
	InstHld             interface{} `json:"inst_hld"`
	Pb                  float64     `json:"pb"`
	PeTtm               float64     `json:"pe_ttm"`
	ContractSize        int         `json:"contract_size"`
	VariableTickSize    string      `json:"variable_tick_size"`
	Time                int64       `json:"time"`
	Open                float64     `json:"open"`
}

type Data struct {
	Quote Quote `json:"quote"`
}
