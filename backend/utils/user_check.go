package utils

import "main/global"

// 检查用户类型是否在允许的类型列表中
func CheckUserType(usertype int, allowedtype ...int) bool {
	for _, allowed := range allowedtype {
		if usertype == allowed {
			return true
		}
	}
	return false
}

// 检测账户是否存在
func IsExists(model interface{}, field string, value interface{}) (bool, error) {
	var count int64

	if err := global.DB.Model(model).Where(field+" = ?", value).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}
