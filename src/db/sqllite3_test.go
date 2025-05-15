package db

import (
	"testing"
)

func TestInitDb(t *testing.T) {
	var result Sto
	dbLite.Table("stock_basic").Where("1=1").Limit(1).Scan(&result)
	t.Log(result)
}
