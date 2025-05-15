package stock

import (
	"cat_ben/src/db"
	"cat_ben/src/pkg/bizcall"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

// from https://xueqiu.com/S/INTC

var sleep = 0

func getSleep() int {
	if sleep > 0 {
		return sleep
	}
	value, err := db.GetValue(`sleep`)
	if err != nil {
		sleep = 1000
	}
	atoi, err := strconv.Atoi(value)
	if err != nil {
		sleep = 1000
	}
	sleep = atoi
	log.Printf("sleep:%d", sleep)
	return sleep
}

func FlushBasic(hard string, tpe string) {
	log.Println("flush start")
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover:", r)
		}
	}()
	stocks, _ := db.GetAllStockFromDB(hard, tpe)
	begin := time.Now().Unix()
	count := 0
	for _, stock := range stocks {
		time.Sleep(time.Duration(getSleep()) * time.Millisecond)
		FlushOne(stock)
		count++
	}

	log.Print("flush finish cost:")
	log.Println(time.Now().Unix() - begin)
	log.Print("count:")
	log.Println(count)
}

func FlushOne(stock *db.Sto) {
	pe, yield, chn, price, h52, l52, liangbi, shizhi, huanshoulv, err := getDetailFromXQ(stock.Name)
	if err != nil {
		return
	}
	hl := (price - l52) / (h52 - l52)
	if h52 == l52 {
		hl = 1
	}
	hl = math.Round(hl*10000) / 10000
	kData, cjl, _ := getKlineFromXQ(stock.Name, "day", 69)
	if len(kData) > 0 {
		ts := kData[len(kData)-1].Ts / 1000
		if time.Now().Unix()-ts > 86400*30 && ts > 0 {
			err := db.DeleteStoById(stock.ID)
			if err != nil {
				log.Printf("delete one stock err,%s:%s", stock.Name, err)
			}
			log.Println("delete one stock:" + stock.Name)
			return
		}
	}

	zcl := getZhicheng(kData, 0.009)
	kDataW, _, _ := getKlineFromXQ(stock.Name, "week", 159)
	zcw := getZhicheng(kDataW, 0.012)

	err = db.UpdateByID(stock.ID, pe, yield, chn, price, h52, l52, hl, liangbi, shizhi, huanshoulv, cjl, zcl, zcw)
	if err != nil {
		log.Println(err)
	}
}

// 这段代码是获取日 周级别k线的代码
// https://stock.xueqiu.com/v5/stock/chart/kline.json?
// symbol=INTC&begin=1712891301705&period=week&type=before&
// count=-284&indicator=kline,pe,pb,ps,pcf,market_capital,agt,ggt,balance
//
//period:week day month
func getKlineFromXQ(name string, period string, count int64) (data []*KlineData, cjl float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover:", r)
			log.Println("name:", name)
		}
	}()
	token, _ := db.GetValue("xueqiu_token")
	header := map[string]string{
		"cookie":     token,
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	}
	url := fmt.Sprintf("https://stock.xueqiu.com/v5/stock/chart/kline.json?symbol=%s&begin=%d&period=%s&type=before&count=-%d&indicator=kline,pe,pb,ps,pcf,market_capital,agt,ggt,balance",
		name, time.Now().UnixMilli(), period, count)
	resp, err := bizcall.GetJSONWithHeader(context.Background(), url, header)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}
	var respXQ KlineJSONDataXQ
	json.Unmarshal(resp, &respXQ)
	items := respXQ.Data.Item
	length := len(items)
	data = make([]*KlineData, length)
	vtot := float64(0)
	pe := float64(0)
	for i, v := range items {
		ts := int64(v[0])
		if v[12] != float64(0) {
			pe = v[12]
		}
		d := &KlineData{
			Ts:      ts,
			Volume:  v[1],
			Open:    v[2],
			High:    v[3],
			Low:     v[4],
			Close:   v[5],
			Chg:     v[6],
			Percent: v[7],
			Pe:      pe,
			Time:    time.UnixMilli(ts).Format("2006-01-02 15-04-05"),
		}
		data[i] = d
		vtot += d.Volume
	}
	vavg := vtot / float64(length)
	vcurrent := items[length-1][1]
	if len(items) > 2 {
		vcurrent = (items[length-1][1] + items[length-2][1] + items[length-3][1]) / 3
	}
	cjl = float64(-1)
	if vavg != 0 {
		cjl = vcurrent / vavg
	}
	return data, cjl, nil
}

// timestamp", "volume", "open", "high", "low", "close",
// "chg", "percent", "turnoverrate", "amount", "volume_post", "amount_post
type KlineData struct {
	Ts      int64   `json:"timestamp"`
	Volume  float64 `json:"volume"`
	Open    float64 `json:"open"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Close   float64 `json:"close"`
	Chg     float64 `json:"chg"`     //涨额
	Percent float64 `json:"percent"` //涨幅
	Pe      float64 `json:"pe"`      //涨幅
	Time    string
}
type KlineJSONDataXQ struct {
	Data             KlineDataXQ `json:"data"`
	ErrorCode        int         `json:"error_code"`
	ErrorDescription string      `json:"error_description"`
}
type KlineDataXQ struct {
	Symbol string      `json:"symbol"`
	Column []string    `json:"column"`
	Item   [][]float64 `json:"item"`
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
	_ = json.Unmarshal(resp, &respXQ)
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
	if price == 0 {
		return 0, 0, "", 0, 0, 0, 0, 0, 0, errors.New("price err")
	}
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
