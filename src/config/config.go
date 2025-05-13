package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type jsonConfig struct {
	DbPath string `json:"dbPath"`
	Port   int    `json:"port"`
	Debug  bool   `json:"debug"`
}

var Config jsonConfig

func LoadAll() {
	loadLog()
	loadConfig()
}

func loadConfig() {
	cwd, _ := os.Getwd()
	file, err := os.ReadFile(filepath.Join(cwd, "config.json"))
	if err != nil {
		file, err = os.ReadFile(filepath.Join(cwd, "../../config.json"))
		if err != nil {
			panic("config.json missing")
		}
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		panic("config.json format err")
	}
	log.Printf("load config success:%+v", Config)
}

func loadLog() {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "app.log")
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("loadLog err")
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
