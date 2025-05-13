package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	LoadAll()
	t.Log(Config)
}
