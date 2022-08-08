package bizcall

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpCliWithTimeout = &http.Client{Timeout: time.Second * 30}

func PostFormWithHeader(ctx context.Context, url string, info url.Values, header map[string]string) ([]byte, error) {
	log.Printf("post form %s req:%s, header:%s", url, info.Encode(), header)
	req, _ := http.NewRequest("POST", url, strings.NewReader(info.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("post form call %s err:%s", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	log.Printf("post form %s req:%s, header:%s, resp:%s", url, info.Encode(), header, string(respByte))
	return respByte, nil
}

func PostJSONWithHeader(ctx context.Context, url string, info interface{}, header map[string]string) ([]byte, error) {
	msg, _ := json.Marshal(info)
	log.Printf("post json %s req:%s, header:%s", url, string(msg), header)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(msg))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := httpCliWithTimeout.Do(req)
	if err != nil {
		log.Printf("post json %s err:%s", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	//logs.CtxInfo(ctx, "post json %s req:%s, header:%s, resp:%s", url, string(msg), header, string(respByte))
	return respByte, nil
}

func GetJSONWithHeader(ctx context.Context, url string, header map[string]string) ([]byte, error) {
	log.Printf("get json %s header:%s", url, header)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := httpCliWithTimeout.Do(req)
	if err != nil {
		log.Printf("get json %s err:%s", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	log.Printf("get json %s, header:%s, resp:%s", url, header, string(respByte))
	return respByte, nil
}
