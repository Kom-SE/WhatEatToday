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

// 验证密码
func CheckPW(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 定义双token的过期时间、签发者、密钥
const (
	ATokenExpiredDuration = 2 * time.Hour
	RTokenExpiredDuration = 30 * 24 * time.Hour
	TokenIssuer           = "admin"
)

// 定义 JWT 令牌的声明
var (
	TokenSecret  = []byte("secret for gods")
	ErrTokenTime = errors.New("token time error")
)

// 载荷声明
type PayLoad struct {
	UserID   uint  `json:"user_id"`
	Usertype uint8 `json:"usertype"`
	jwt.RegisteredClaims
}

// 检查 JWT 令牌的签名方法
func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return []byte("secret for gods"), nil
}

// 生成JWT令牌
func GenerateJWT(userid uint, usertype uint8) (string, string, error) {
	// 生成访问令牌
	// 构建 凭证 基础信息
	rc := jwt.RegisteredClaims{
		Issuer:    TokenIssuer,                                               // 颁发人
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ATokenExpiredDuration)), // 到期时间
	}
	//绑定载荷信息
	atClaims := PayLoad{
		UserID:           userid,
		Usertype:         usertype,
		RegisteredClaims: rc,
	}
	// 使用SHA256对载荷非对称加密，进行签名和加盐
	atoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString(TokenSecret)
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	rt := rc
	// 设置刷新令牌的过期时间
	rt.ExpiresAt = jwt.NewNumericDate(time.Now().Add(RTokenExpiredDuration))
	// 绑定载荷信息
	rtClaims := PayLoad{
		UserID:           userid,
		Usertype:         usertype,
		RegisteredClaims: rt,
	}
	// 生成刷新令牌
	rtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims).SignedString(TokenSecret)
	if err != nil {
		return "", "", err
	}

	// 返回访问令牌和刷新令牌
	return "Bearer " + atoken, "Bearer " + rtoken, nil
}

// 验证JWT令牌
func ParseJWT(tokenString string) (uint8, uint8, error) {
	var tPayLoad PayLoad
	// 解析令牌
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	token, err := jwt.ParseWithClaims(tokenString, &tPayLoad, keyFunc)
	if err != nil {
		return 0, 3, errors.New("invalid token")
	}
	// 验证令牌是否有效
	if !token.Valid {
		return 0, 3, errors.New("invalid token")
	}
	// 检查过期时间
	if tPayLoad.ExpiresAt != nil && tPayLoad.ExpiresAt.Time.Before(time.Now()) {
		return 0, 3, ErrTokenTime
	}
	// 验证令牌不是在未来颁发的
	if tPayLoad.IssuedAt != nil && tPayLoad.IssuedAt.After(time.Now()) {
		return 0, 3, ErrTokenTime
	}

	return uint8(tPayLoad.UserID), tPayLoad.Usertype, nil
}

// 根据rtoken刷新atoken
func RefreshAToken(rtoken string) (string, string, error) {
	// 去除Bearer前缀
	if len(rtoken) > 7 && rtoken[:7] == "Bearer " {
		rtoken = rtoken[7:]
	}

	// 解析刷新令牌rtoken
	userid, usertype, err := ParseJWT(rtoken)
	if err != nil {
		return "", "", err
	}

	// 生成新的访问令牌和刷新令牌
	atoken, rtoken, err := GenerateJWT(uint(userid), usertype)
	if err != nil {
		return "", "", err
	}

	return "Bearer " + atoken, "Bearer " + rtoken, nil
}
