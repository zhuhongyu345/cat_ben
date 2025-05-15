package option

import (
	"encoding/json"
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
				_ = json.NewEncoder(w).Encode(`{}`)
				return
			}
			chain := GetOptionChain(atoi)
			log.Printf("resp:%s", id)
			w.Header().Add("Access-Control-Allow-Origin", "*")
			_ = json.NewEncoder(w).Encode(chain)
		}),
	}
	_ = server.ListenAndServe()
}

func TestGetOptionChain(t *testing.T) {
	chain := GetOptionChain(913243251)
	marshal, _ := json.Marshal(chain)
	t.Log(string(marshal))
}

func TestDataFromUws(t *testing.T) {
	DataFromUws()
}
