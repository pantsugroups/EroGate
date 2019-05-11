package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func ManualLogin(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		return c.String(http.StatusInternalServerError, "parse error.")
	}
	u, err := LoginRequests(c, pathMap[c.Path()], c.Request().Form.Encode())
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
	if session, ok := c.Request().Header["X-Headers-Session"]; ok {
		u, err := ParseToken(session[0])
		if err != nil {
			return c.String(http.StatusOK, "auth error.")
		}
		err = c.Request().ParseForm()
		if err != nil {
			return c.String(http.StatusInternalServerError, "parse error.")
		}
		form := c.Request().Form.Encode()
		err = ServerRequests(c, pathMap[c.Path()], *u, form)
		if err != nil {
			return c.String(http.StatusInternalServerError, "InternalServerError.")
		}

		return nil
	} else {

		// 返回403
		return c.String(http.StatusForbidden, "Please.Login!")
	}

}
