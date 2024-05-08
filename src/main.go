package main

import (
	"cat_ben/src/db"
	"cat_ben/src/option"
	"cat_ben/src/stock"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/chainEcharts", http.HandlerFunc(OptionServer))
	mux.Handle("/search", http.HandlerFunc(SelectServer))
	mux.Handle("/history", http.HandlerFunc(HistoryServer))
	mux.Handle("/flush", http.HandlerFunc(FlushServer))
	mux.Handle("/deleteOne", http.HandlerFunc(DeleteServer))
	mux.Handle("/config", http.HandlerFunc(ConfigServer))
	http.ListenAndServe(":8001", mux)
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	log.Printf("delete req id:" + id)
	i, _ := strconv.ParseInt(id, 10, 64)
	_ = db.DeleteStoById(i)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, string(`success`))
	return
}

func ConfigServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConfigServer")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	respData := `success`

	rd := r.FormValue("rd")
	key := r.FormValue("key")
	value := r.FormValue("value")
	if rd == "read" {
		resp, _ := db.GetValue(key)
		respData = fmt.Sprintf(`{"key":"%s","value":"%s"}`, key, resp)
	}
	if rd == "write" {
		_ = db.UpdateValue(key, value)
		respData = fmt.Sprintf(`{"key":"%s","value":"%s"}`, key, value)
	}
	log.Printf("ConfigServer:" + respData)
	fmt.Fprintf(w, respData)
	return
}

func FlushServer(w http.ResponseWriter, r *http.Request) {
	hard := r.FormValue("hard")
	tpe := r.FormValue("type")
	log.Printf("flush req hard:" + hard)
	go stock.FlushBasic(hard, tpe)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, string(`success`))
	return
}

func HistoryServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("HistoryServer")
	name := r.FormValue("name")
	count := r.FormValue("count")
	period := r.FormValue("period")
	resp := stock.GetHistory(strings.ToUpper(name), period, count)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	marshal, err := json.Marshal(resp)
	if err != nil {
		log.Printf("%s", err)
	}
	fmt.Fprintf(w, string(marshal))
	return
}
func SelectServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("SelectServer")
	zclHigh := r.FormValue("zclHigh")
	zclLow := r.FormValue("zclLow")
	cjlLow := r.FormValue("cjlLow")
	cjlHigh := r.FormValue("cjlHigh")
	hlLow := r.FormValue("hlLow")
	hlHigh := r.FormValue("hlHigh")
	peHigh := r.FormValue("peHigh")
	peLow := r.FormValue("peLow")
	name := r.FormValue("name")
	yield := r.FormValue("yield")
	priceLow := r.FormValue("priceLow")
	priceHigh := r.FormValue("priceHigh")
	liangbi := r.FormValue("liangbi")
	tpe := r.FormValue("type")
	skip := r.FormValue("skip")
	size := r.FormValue("size")
	sort := r.FormValue("sort")
	sortType := r.FormValue("sortType")
	search := stock.Search(name, zclLow, zclHigh, cjlLow, cjlHigh, hlLow, hlHigh, peHigh, peLow, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	marshal, err := json.Marshal(search)
	log.Printf("%s", err)
	fmt.Fprintf(w, string(marshal))
	return
}

func OptionServer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	atoi, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, `{}`)
		return
	}
	chain := option.GetOptionChain(atoi)
	marshal, _ := json.Marshal(chain)
	log.Printf("resp:%s", id)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, string(marshal))
	return
}
