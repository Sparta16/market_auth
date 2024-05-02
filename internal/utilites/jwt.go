package utilites

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	models "weather_back_api_getway/internal"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.UUID
	claims["login"] = user.Login
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte("1231231321")) //TODO: в конфиге добавить секрет
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
