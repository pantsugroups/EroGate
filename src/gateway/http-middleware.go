package main

import (
	"github.com/labstack/echo"
	"net/http"
)
func ManualLogin(c echo.Context) error{
	return c.String(http.StatusOK,"this is a special for verify the login event.")
}
func ManualGateWay(c echo.Context) error {
	if session, ok := c.Request().Header["x-headers-session"]; ok {
		u,err := ParseToken(session[0])
		if err != nil{
			return c.String(http.StatusOK,"auth error.")
		}
		return c.String(http.StatusOK,"welcome:"+u.Username)
	} else {
		// 返回403
		return c.String(http.StatusForbidden, "Please.Login!")
	}
	return nil //明明理論上執行不到這裏，卻還是給我報提示mmm
}
