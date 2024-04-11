package db

type ConfigPro struct {
	Key   string `gorm:"column:k"`
	Value string `gorm:"column:v"`
}

func GetValue(key string) (string, error) {
	if dbLite == nil {
		InitDb()
	}
	resp := make([]*ConfigPro, 0)
	err := dbLite.Table("config_pro").Where("k=?", key).Find(&resp).
		Order("id desc").Error
	if len(resp) > 0 {
		return resp[0].Value, nil
	}
	return "", err
}

func UpdateValue(key, value string) error {
	if dbLite == nil {
		InitDb()
	}
	update := map[string]interface{}{
		"v": value,
	}
	err := dbLite.Table("config_pro").Where("k = ?", key).Updates(update).Error
	return err

}
