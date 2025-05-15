package server

import (
	"cat_ben/src/db"
	"cat_ben/src/option"
	"cat_ben/src/stock"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AddServer(w http.ResponseWriter, r *http.Request) {
	name := strings.ToUpper(r.FormValue("name"))
	tpe := r.FormValue("type")
	i, _ := strconv.Atoi(tpe)
	sto := &db.Sto{
		Name: name,
		Type: i,
		TAG:  1,
	}
	nyse := stock.GetFromNyse(name)
	if len(nyse) > 0 {
		sto.Mic = nyse[0].MicCode
	}
	byName, _ := db.SelectStoByName(name)
	if byName != nil && byName.Name != "" {
		sto = byName
	} else {
		_ = db.CreateStos([]*db.Sto{sto})
	}
	stock.FlushOne(sto)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode("success")
	return
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	log.Printf("delete req id:" + id)
	i, _ := strconv.ParseInt(id, 10, 64)
	_ = db.DeleteStoById(i)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode("success")
	return
}
func TagServer(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	tag := r.FormValue("tag")
	log.Printf("tag req id:" + id)
	i, _ := strconv.ParseInt(id, 10, 64)
	_ = db.UpdateTagByID(i, tag)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode("success")
}

func ConfigQueryServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConfigServer")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	kv, err := db.SelectAllKV()
	log.Printf("ConfigServer:%v", kv)
	if err != nil {
		_ = json.NewEncoder(w).Encode(`[]`)
	}
	_ = json.NewEncoder(w).Encode(kv)

}
func ConfigUpdateServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConfigServer")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	key := r.FormValue("key")
	val := r.FormValue("val")
	_ = db.UpdateValue(key, val)
	_ = json.NewEncoder(w).Encode("success")
}

func HistoryServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("HistoryServer")
	name := r.FormValue("name")
	count := r.FormValue("count")
	period := r.FormValue("period")
	resp := stock.GetHistory(strings.ToUpper(name), period, count)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)

}
func SelectServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("SelectServer")
	param := &db.SearchDto{
		ZclHigh:   getFloat(r.FormValue("zclHigh")),
		ZclLow:    getFloat(r.FormValue("zclLow")),
		CjlLow:    getFloat(r.FormValue("cjlLow")),
		CjlHigh:   getFloat(r.FormValue("cjlHigh")),
		HlLow:     getFloat(r.FormValue("hlLow")),
		HlHigh:    getFloat(r.FormValue("hlHigh")),
		PeHigh:    getFloat(r.FormValue("peHigh")),
		Yield:     getFloat(r.FormValue("yield")),
		PriceLow:  getFloat(r.FormValue("priceLow")),
		PriceHigh: getFloat(r.FormValue("priceHigh")),
		Liangbi:   getFloat(r.FormValue("liangbi")),
		Name:      r.FormValue("name"),
		Tpe:       getInt(r.FormValue("type")),
		Skip:      getInt(r.FormValue("skip")),
		Size:      getInt(r.FormValue("size")),
		SortType:  r.FormValue("sortType"),
		Sort:      r.FormValue("sort"),
	}
	pel := r.FormValue("peLow")
	if pel == "" {
		param.PeLow = -10000
	} else {
		param.PeLow = getFloat(r.FormValue("peLow"))
	}

	search := stock.Search(param)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(search)
}

func OptionServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("id")
	atoi, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = json.NewEncoder(w).Encode(`{}`)
	}
	chain := option.GetOptionChain(atoi)
	_ = json.NewEncoder(w).Encode(chain)
}

func FlushServer(w http.ResponseWriter, r *http.Request) {
	hard := r.FormValue("hard")
	tpe := r.FormValue("type")
	log.Printf("flush req hard:" + hard)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	go stock.FlushBasic(hard, tpe)
	_ = json.NewEncoder(w).Encode(stock.Doing)
}
