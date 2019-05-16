package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"

	"io/ioutil"
	"log"
	"net/http"
	"time"
)
import "github.com/labstack/echo"
import "github.com/labstack/echo/middleware"

func ServerRequests(c echo.Context, backend string, user UserInfo, form string) error {

	body_, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	encodeString := base64.StdEncoding.EncodeToString(body_)

	data, err := json.Marshal(c.Request().Header)
	if err != nil {
		return err
	}

	r := &ForwardRequest{UserInfo: user, RequestHeader: string(data), RequestBody: encodeString, Method: c.Request().Method, RequestForm: form}
	bytesData, err := json.Marshal(r)
	if err != nil {
		return err
	}
	backend = ParseUrl(backend)
	res, err := http.Post(backend+c.Request().URL.Path, "application/json;charset=utf-8", bytes.NewReader(bytesData))
	if err != nil {
		return err
	}

	defer func() {
		if err = c.Request().Body.Close(); err != nil {
			//log
		} else if err = res.Body.Close(); err != nil {
			//log
		} else {
			err = nil
		}

	}()
	body_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	for i, j := range res.Header {
		if len(j) != 1 {
			for _, l := range j {
				c.Response().Header().Add(i, l)
			}
		} else {
			c.Response().Header().Add(i, j[0])
		}

	}
	c.Response().WriteHeader(res.StatusCode)
	_, err = c.Response().Write(body_)
	if err != nil {
		return err
	}
	return err
}

//func LoginRequests(c echo.Context, backend string, form string) (UserInfo, error) {
//	body_, err := ioutil.ReadAll(c.Request().Body)
//	if err != nil {
//		return UserInfo{}, err
//	}
//	encodeString := base64.StdEncoding.EncodeToString(body_)
//
//	data, err := json.Marshal(c.Request().Header)
//	if err != nil {
//		return UserInfo{}, err
//	}
//
//	r := &ForwardRequest{RequestHeader: string(data), RequestBody: encodeString, Method: c.Request().Method, RequestForm: form}
//	bytesData, err := json.Marshal(r)
//	if err != nil {
//		return UserInfo{}, err
//	}
//	backend = ParseUrl(backend)
//	res, err := http.Post(backend+c.Request().URL.Path, "application/json;charset=utf-8", bytes.NewReader(bytesData))
//	if err != nil {
//		return UserInfo{}, err
//	}
//
//	defer func() {
//		if err = c.Request().Body.Close(); err != nil {
//			//log
//		} else if err = res.Body.Close(); err != nil {
//			//log
//		} else {
//			err = nil
//		}
//
//	}()
//	body_, err = ioutil.ReadAll(res.Body)
//	if err != nil {
//		return UserInfo{}, err
//	}
//	var v Verify
//	fmt.Println(string(body_))
//	err = json.Unmarshal(body_, &v)
//	if err != nil {
//		return UserInfo{}, err
//	}
//	if v.Secret == conf.Base.Secret {
//		return v.U, nil
//	}
//
//	return UserInfo{}, nil
//}

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
		err = errors.New("cannot convert claim to map claim")
		return user, err
	}
	if !token.Valid {
		err = errors.New("token is invalid")
		return user, err
	}

	user.ID = int(claim["id"].(float64))
	user.Username = claim["username"].(string)
	return user, nil
}

func StartEchoHandle() {
	// Echo instance

	e = echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	err := LoadRoutes()
	if err != nil {
		log.Println("Load routes error", err)
	}

	// Start server
	e.Logger.Fatal(e.Start(":" + conf.Base.Port))
}
