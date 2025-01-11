package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("123")

func GenerateToken(userID int) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),                               // Chuyển ID thành chuỗi
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token hết hạn sau 24 giờ
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if ok && token.Valid {
		// Chuyển đổi `claims.Subject` từ `string` sang `int`
		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}

	return 0, jwt.ErrSignatureInvalid
}

// VerifyToken kiểm tra tính hợp lệ của token
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Giải mã và xác thực token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra phương pháp ký mã hóa
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("phương pháp ký không hợp lệ")
		}
		return jwtKey, nil
	})

	// Nếu token không hợp lệ hoặc xảy ra lỗi
	if err != nil {
		return nil, err
	}

	// Kiểm tra token có hợp lệ không
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Kiểm tra token hết hạn
		expiration := int64(claims["exp"].(float64))
		if expiration < time.Now().Unix() {
			return nil, errors.New("token đã hết hạn")
		}
		return token, nil
	}

	return nil, errors.New("token không hợp lệ")
}
