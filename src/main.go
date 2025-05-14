package main

import (
	"cat_ben/src/chromedriver"
	"cat_ben/src/config"
	"cat_ben/src/db"
	"cat_ben/src/option"
	"cat_ben/src/stock"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	config.LoadAll()

	log.Printf("start server")
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(config.Config.Static))))
	mux.Handle("/chainEcharts", http.HandlerFunc(OptionServer))
	mux.Handle("/search", http.HandlerFunc(SelectServer))
	mux.Handle("/history", http.HandlerFunc(HistoryServer))
	mux.Handle("/flush", http.HandlerFunc(FlushServer))
	mux.Handle("/deleteOne", http.HandlerFunc(DeleteServer))
	mux.Handle("/addOne", http.HandlerFunc(AddServer))
	mux.Handle("/tagOne", http.HandlerFunc(TagServer))
	mux.Handle("/config", http.HandlerFunc(ConfigServer))
	log.Printf("start listen:%d", config.Config.Port)
	go FlushTask()
	err := http.ListenAndServe(":"+strconv.Itoa(config.Config.Port), mux)
	if err != nil {
		log.Printf("start err:%v", err)
	}
}

var doing = false

func FlushTask() {
	for {
		time.Sleep(time.Hour * 1)
		doing = true
		chromedriver.GetTokenAndSave()
		stock.FlushBasic("1", "")
		doing = false
	}
}

func AddServer(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	tpe := r.FormValue("type")
	i, _ := strconv.Atoi(tpe)
	stos := make([]*db.Sto, 0)
	stos = append(stos, &db.Sto{
		Name: name,
		Type: i,
		TAG:  1,
	})
	_ = db.CreateStos(stos)
	stock.FlushBasic("1", "-1")
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
	_ = json.NewEncoder(w).Encode(respData)

}

func FlushServer(w http.ResponseWriter, r *http.Request) {
	hard := r.FormValue("hard")
	tpe := r.FormValue("type")
	log.Printf("flush req hard:" + hard)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if !doing || tpe == "-1" {
		go stock.FlushBasic(hard, tpe)
		_ = json.NewEncoder(w).Encode("success")
	} else {
		_ = json.NewEncoder(w).Encode("doing")
	}
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
