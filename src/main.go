package main

import (
	"cat_ben/src/config"
	"cat_ben/src/server"
	"log"
	"net/http"
	"strconv"
)

func main() {
	//task
	go server.FlushTask()

	//server
	log.Printf("start server")
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static",
		http.FileServer(http.Dir(config.Config.Static))))
	mux.Handle("/chainEcharts", http.HandlerFunc(server.OptionServer))
	mux.Handle("/search", http.HandlerFunc(server.SelectServer))
	mux.Handle("/history", http.HandlerFunc(server.HistoryServer))
	mux.Handle("/flush", http.HandlerFunc(server.FlushServer))
	mux.Handle("/deleteOne", http.HandlerFunc(server.DeleteServer))
	mux.Handle("/addOne", http.HandlerFunc(server.AddServer))
	mux.Handle("/tagOne", http.HandlerFunc(server.TagServer))
	mux.Handle("/configQuery", http.HandlerFunc(server.ConfigQueryServer))
	mux.Handle("/configUpdate", http.HandlerFunc(server.ConfigUpdateServer))
	log.Printf("start listen:%d", config.Config.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(config.Config.Port), mux)
	if err != nil {
		log.Printf("start err:%v", err)
	}
}
