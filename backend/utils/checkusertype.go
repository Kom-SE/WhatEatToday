package utils

func CheckUserType(usertype int, allowedtype ...int) bool {
	for _, allowed := range allowedtype {
		if usertype == allowed {
			return true
		}
	}
	return false
}
