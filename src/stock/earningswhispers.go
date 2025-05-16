package stock

import (
	"cat_ben/src/config"
	"cat_ben/src/pkg/bizcall"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SaveImgLocal(code string) {
	filePath := config.Config.Static + "/icon/" + strings.ToUpper(code) + ".png"
	_, err := os.Stat(filePath)
	if err == nil {
		return
	}
	url := "https://www.earningswhispers.com/api/icon/" + strings.ToUpper(code)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Referer", "https://www.earningswhispers.com")
	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()

	// 将响应内容写入文件
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	out, err := os.Create(filePath)
	if err != nil {
		return
	}
	log.Printf("%s:%s", code, filePath)
	_, err = io.Copy(out, res.Body)
}

// GetByDate 20250519
func GetByDate(date string) ([]*EWResp, error) {
	url := "https://www.earningswhispers.com/api/caldata/" + date
	resp, err := bizcall.GetJSONWithHeader(context.Background(), url, map[string]string{
		"referer": "https://www.earningswhispers.com/calendar/20250519",
	})
	if err != nil {
		log.Printf("earningswhispers call err:%s", err)
		return nil, err
	}
	var ew []*EWResp
	err = json.Unmarshal(resp, &ew)
	return ew, err
}

type EWResp struct {
	Ticker      string  `json:"ticker"`
	Company     string  `json:"company"`
	Total       int     `json:"total"`
	NextEPSDate string  `json:"nextEPSDate"`
	ReleaseTime int     `json:"releaseTime"`
	QDate       string  `json:"qDate"`
	Q1RevEst    float64 `json:"q1RevEst"`
	Q1EstEPS    float64 `json:"q1EstEPS"`
	ConfirmDate string  `json:"confirmDate"`
	EpsTime     string  `json:"epsTime"`
	QuarterDate string  `json:"quarterDate"`
	QSales      float64 `json:"qSales"`
}
