package utils

import "main/global"

// 检测账户是否存在
func IsExists(model interface{}, field string, value interface{}) (bool, error) {
	var count int64

	if err := global.DB.Model(model).Where(field+" = ?", value).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}
