package option

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func TestGetOptionChainHttp(t *testing.T) {

	server := http.Server{
		Addr: ":8001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Query().Get("id")
			atoi, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				fmt.Fprintf(w, `{}`)
				return
			}
			chain := GetOptionChain(atoi)
			marshal, _ := json.Marshal(chain)
			log.Printf("resp:%s", id)
			w.Header().Add("Access-Control-Allow-Origin", "*")
			fmt.Fprintf(w, string(marshal))
			return
		}),
	}
	server.ListenAndServe()
}

func TestGetOptionChain(t *testing.T) {
	chain := GetOptionChain(913243251)
	marshal, _ := json.Marshal(chain)
	t.Log(string(marshal))
}

func TestDataFromUws(t *testing.T) {
	DataFromUws()
}
