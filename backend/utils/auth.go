package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 密码bcrypt加密
func BcryptPW(password string) (string, error) {
	hashpw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hashpw), err
}

// 生成JWT令牌
func GenerateJWT(username string, usertype uint8) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"usertype": usertype,
		"exp":      time.Now().Add((time.Hour * 72)).Unix()})

	sigenedtoken, err := token.SignedString([]byte("secret for gods"))
	return "Bearer " + sigenedtoken, err
}

// 验证密码
func CheckPW(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 验证JWT令牌
func ParseJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("secret for gods"), nil
	})

	if err != nil {
		return "", nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return username, nil
	}
	return "", errors.New("invalid token")
}
