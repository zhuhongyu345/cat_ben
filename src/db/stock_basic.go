package db

type Sto struct {
	ID    int64   `gorm:"column:id"`
	Name  string  `gorm:"column:name"`
	Type  int     `gorm:"column:type"`
	CHN   string  `gorm:"column:chn"`
	Yield float64 `gorm:"column:yield"`
	PE    float64 `gorm:"column:pe"`
	Price float64 `gorm:"column:price"`
	H52   float64 `gorm:"column:h52"`
	L52   float64 `gorm:"column:l52"`
	Hl    float64 `gorm:"column:hl"`
}

func CreateStos(stos []*Sto) error {
	if dbLite == nil {
		InitDb()
	}
	err := dbLite.Table("stock_basic").Create(stos).Error
	return err
}

func GetAllStockFromDB() ([]*Sto, error) {
	if dbLite == nil {
		InitDb()
	}
	resp := make([]*Sto, 0)
	err := dbLite.Table("stock_basic").Where("1=1 and (hl is null or shizhi is null)").Find(&resp).Error
	return resp, err
}

func UpdateByID(id int64, pe, yield float64, chn string, price, h52, l52, hl, liangbi, shizhi, huanshoulv float64) error {
	if dbLite == nil {
		InitDb()
	}
	update := map[string]interface{}{
		"pe":         pe,
		"yield":      yield,
		"chn":        chn,
		"price":      price,
		"h52":        h52,
		"l52":        l52,
		"hl":         hl,
		"liangbi":    liangbi,
		"shizhi":     shizhi,
		"huanshoulv": huanshoulv,
	}
	err := dbLite.Table("stock_basic").Where("id = ?", id).Updates(update).Error
	return err

}
