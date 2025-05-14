package chromedriver

import (
	"cat_ben/src/config"
	"testing"
)

func TestGetTokenAndSave(t *testing.T) {
	config.LoadAll()
	GetTokenAndSave()

}
