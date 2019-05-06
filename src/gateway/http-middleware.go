package gateway

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Base.Secret), nil
	}
}

func CreateToken(user *UserInfo) (tokens string, err error) {

	claim := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokens, err = token.SignedString([]byte(conf.Base.Secret))
	if err != nil {
		return "", nil
	}
	return tokens, nil
}

func ParseToken(tokens string) (user *UserInfo, err error) {
	user = &UserInfo{}
	token, err := jwt.Parse(tokens, secret())
	if err != nil {
		return user, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return user, err
	}
	if !token.Valid {
		err = errors.New("token is invalid")
		return user, err
	}

	user.ID = claim["id"].(int)
	user.Username = claim["username"].(string)
	return user, nil
}
func Verify(session string) {

}

func ManualGateWay(c echo.Context) error {

	if _, ok := c.Request().Header["x-headers-session"]; ok {
		// 驗證session
	} else {
		// 返回403
		return c.String(http.StatusForbidden, "Please.Login!")
	}
	return nil //明明理論上執行不到這裏，卻還是給我報提示mmm
}
