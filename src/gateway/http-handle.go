package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
)

func ManualLogin(c echo.Context) error {
	var v Verify

	body_, err := ioutil.ReadAll(c.Request().Body)
	defer func() {
		err = c.Request().Body.Close()
		if err != nil {
			// handle err .
			log.Println(err) // error level is not hard.
		}
	}()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body_, &v)
	if err != nil {
		return err
	}
	if v.Secret == conf.Base.Secret {
		token, err := CreateToken(&v.U)
		if err != nil {
			return nil
		}
		return c.String(http.StatusOK, token)
	}
	return c.String(http.StatusOK, "")
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
