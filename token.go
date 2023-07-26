package ginz

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"

	"github.com/qf0129/ginz/pkg/encrypt"
)

var ErrTokenExpired = errors.New("TokenExpired")

// 创建令牌
func CreateToken(body string) (string, error) {
	tokenByte := make([]byte, len(body)+8) // 前8位放时间戳
	binary.BigEndian.PutUint64(tokenByte[:8], uint64(time.Now().Unix()))
	for i, v := range body {
		tokenByte[i+8] = byte(v)
	}

	token, err := encrypt.DesEncrypt(tokenByte, []byte(Config.SecretKey))
	return hex.EncodeToString(token), err
}

// 解析令牌
func ParseToken(token string) (string, error) {
	text, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}
	tokenStr, err := encrypt.DesDecrypt([]byte(text), []byte(Config.SecretKey))
	if err != nil {
		return "", err
	}

	tokenTs := int64(binary.BigEndian.Uint64(tokenStr[0:8]))
	nowTs := time.Now().Unix()

	if nowTs > tokenTs+int64(Config.TokenExpiredTime) {
		return "", ErrTokenExpired
	}

	body := string(tokenStr[8:])
	return body, nil
}
