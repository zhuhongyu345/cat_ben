package db

type ConfigProgram struct {
	Key   string `gorm:"column:k"`
	Value string `gorm:"column:v"`
}

func GetValue(key string) (string, error) {
	resp := make([]*ConfigProgram, 0)
	err := dbLite.Table("config_program").Where("k=?", key).Find(&resp).
		Order("id desc").Error
	if len(resp) > 0 {
		return resp[0].Value, nil
	}
	return "", err
}

func UpdateValue(key, value string) error {
	update := map[string]interface{}{
		"v": value,
	}
	err := dbLite.Table("config_program").Where("k = ?", key).Updates(update).Error
	return err

}

func SelectAllKV() ([]*ConfigProgram, error) {
	resp := make([]*ConfigProgram, 0)
	err := dbLite.Table("config_program").Where("1=1").Find(&resp).Error
	return resp, err
}
