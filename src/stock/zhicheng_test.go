package stock

import (
	"fmt"
	"testing"
)

func TestZhichengwei(t *testing.T) {
	t.Log("test zhicheng")
	kData, _, _ := getKlineFromXQ("KO", "day", 69)
	//去掉几天
	tempKdata := make([]*KlineData, 0)
	for i, v := range kData {
		if i < len(kData)-0 {
			tempKdata = append(tempKdata, v)
		}
	}
	kData = tempKdata
	t.Log(kData[0].Time + " ~ " + kData[len(kData)-1].Time + ": " + fmt.Sprintf(`%d`, len(kData)))

	t.Log(getZhicheng(kData, 0.009))

}
