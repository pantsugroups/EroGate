package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)
import "github.com/labstack/echo"
import "github.com/labstack/echo/middleware"

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
	e.Any(conf.Base.Login,ManualLogin)
	// Start server
	e.Logger.Fatal(e.Start(":" + conf.Base.Port))
}
