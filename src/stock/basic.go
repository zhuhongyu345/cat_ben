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

//from https://xueqiu.com/S/INTC
func flushBasic() {
	stocks, _ := db.GetAllStockFromDB()
	for _, stock := range stocks {
		pe, yield, chn, price, h52, l52, liangbi, shizhi, huanshoulv := getDetailFromXQ(stock.Name)
		hl := (price - l52) / (h52 - l52)
		if h52 == l52 {
			hl = 1
		}
		hl = math.Round(hl*10000) / 10000
		err := db.UpdateByID(stock.ID, pe, yield, chn, price, h52, l52, hl, liangbi, shizhi, huanshoulv)
		if err != nil {
			log.Println(err)
		}
	}
}

func getDetailFromXQ(name string) (float64, float64, string, float64, float64, float64, float64, float64, float64) {
	header := map[string]string{
		"cookie":     "device_id=6be81810d4762245af84672121d4971a; s=dq1nkgp4ed; xq_a_token=bf75ab4bcea18c79de253cb841f2b27e248d8948; xqat=bf75ab4bcea18c79de253cb841f2b27e248d8948; xq_r_token=c7d30dc738a77dd909a8228f3053679e86bf104b; xq_id_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1aWQiOi0xLCJpc3MiOiJ1YyIsImV4cCI6MTY2MTgxNjI0MSwiY3RtIjoxNjU5NDA4MzAyOTA4LCJjaWQiOiJkOWQwbjRBWnVwIn0.YeFdnnx_q_P-m7w2YmzY6oJ7R3eG_260cirgCrFL-tNXH7b-87HEOAzivePXIuw_Vwrw8F0BGvvLoWR687viJO8O6uhbNAMMsnldVEinjDmcsx-eXl9WVUtqi5xrjiPYSyGVRuvwrdVlBn4P2ycqX1_SI9x2IRmfNlA7ZEyuYQxIhYlwJ4K2Af1Lb4Xuo0-iusHEPtvbKHXEcQKsK_adJalgo37FM4MIikkBHtDn1yM_p0CzTjFgeR2V_mgsg1AnC5lI3sXPYMgpp6BvjiAJL46R726gP-roHshjC70253f57zyK8PeLL7YNP65N2zFpTiR9mU_lmedslH4c_IxMYg; u=131659408360534; Hm_lvt_1db88642e346389874251b5a1eded6e3=1659408362; Hm_lpvt_1db88642e346389874251b5a1eded6e3=1659955530",
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	}
	resp, err := bizcall.GetJSONWithHeader(context.Background(), fmt.Sprintf("https://stock.xueqiu.com/v5/stock/quote.json?symbol=%s&extend=detail", name), header)
	log.Println(err)
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
	return ttm, yield, nameCHN, price, h52, l52, liangbi, shizhi, huanshou
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
