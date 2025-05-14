package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

type jsonConfig struct {
	DbPath   string `json:"dbPath"`
	Chrome   string `json:"chromePath"`
	Headless bool   `json:"chromeHeadless"`
	Static   string `json:"staticPath"`
	Port     int    `json:"port"`
	Debug    bool   `json:"debug"`
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
	// 设置日志输出
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "app.log")
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("loadLog err")
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
