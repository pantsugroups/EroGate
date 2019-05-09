package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func ManualLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "" || password == "" {
		return c.String(http.StatusBadRequest, "BadRequest.")
	}

	u, err := LoginRequests(username, password)
	if err != nil {
		return err
	}
	if u == (UserInfo{}) {
		return c.String(http.StatusOK, "Verify Error.")
	}
	token, err := CreateToken(&u)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, token)
}
func ManualGateWay(c echo.Context) error {

	if session, ok := c.Request().Header["x-headers-session"]; ok {
		_, err := ParseToken(session[0])
		if err != nil {
			return c.String(http.StatusOK, "auth error.")
		}

	} else {

		// 返回403
		return c.String(http.StatusForbidden, "Please.Login!")
	}
	return nil
}
