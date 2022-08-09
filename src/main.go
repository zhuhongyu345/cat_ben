package main

import (
	"cat_ben/src/option"
	"cat_ben/src/stock"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/chainEcharts", http.HandlerFunc(OptionServer))
	mux.Handle("/search", http.HandlerFunc(SelectServer))
	http.ListenAndServe(":8001", mux)
}

func SelectServer(w http.ResponseWriter, r *http.Request) {
	hlLow := r.FormValue("hlLow")
	hlHigh := r.FormValue("hlHigh")
	pe := r.FormValue("pe")
	yield := r.FormValue("yield")
	priceLow := r.FormValue("priceLow")
	priceHigh := r.FormValue("priceHigh")
	liangbi := r.FormValue("liangbi")
	tpe := r.FormValue("type")
	skip := r.FormValue("skip")
	size := r.FormValue("size")
	sort := r.FormValue("sort")
	sortType := r.FormValue("sortType")

	log.Printf("req:%s,%s,%s,%s,%s,%s", hlLow, hlHigh, pe, yield, priceLow, priceHigh)
	search := stock.Search(hlLow, hlHigh, pe, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType)
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
