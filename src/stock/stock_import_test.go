package stock

import (
	_ "cat_ben/src/config"
	"testing"
)

func TestBuildStock(t *testing.T) {
	allStockToDB()
}

func TestFlushGuxilv(t *testing.T) {
	//flushBasic()
}
