package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
)

func RegisterHandle(c echo.Context) error {
	var r Register

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
	err = json.Unmarshal(body_, &r)
	if err != nil {
		return err
	}
	if r.Secret == conf.Base.Secret {
		e.Any(r.Route.Route, ManualGateWay)
		authMap[r.Route.Route] = r.Route.Auth
		pathMap[r.Route.Route] = r.Route.BackEnd
		return c.JSON(http.StatusOK, HttpResponse{true, "Load the API Successfully."})
	}
	return c.JSON(http.StatusUnauthorized, HttpResponse{false, "Secret Error"})

}

func ShowAllRoutes(c echo.Context) error {
	//r,err := json.Marshal(pathMap)
	//if err != nil{
	//	return err
	//}
	return c.JSON(http.StatusOK, pathMap)
}

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
		return c.JSON(http.StatusOK, HttpResponse{true, token})
	}
	return c.JSON(http.StatusUnauthorized, HttpResponse{false, "Secret Error"})
}
func ManualGateWay(c echo.Context) error {

	if authMap[c.Path()] == false {

		err := c.Request().ParseForm()
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, HttpResponse{false, "parse error."})
		}
		form := c.Request().Form.Encode()
		err = ServerRequests(c, pathMap[c.Path()], UserInfo{}, form)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, HttpResponse{false, "InternalServerError."})
		}
		return nil
	} else {
		if session, ok := c.Request().Header["X-Headers-Session"]; ok {
			u, err := ParseToken(session[0])
			if err != nil {
				return c.JSON(http.StatusForbidden, HttpResponse{false, "auth error."})
			}
			err = c.Request().ParseForm()
			if err != nil {
				return c.JSON(http.StatusUnprocessableEntity, HttpResponse{false, "parse error."})
			}
			form := c.Request().Form.Encode()
			err = ServerRequests(c, pathMap[c.Path()], *u, form)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, HttpResponse{false, "InternalServerError."})
			}

			return nil
		} else {

			// 返回403
			return c.JSON(http.StatusUnauthorized, HttpResponse{Successful: false, Message: "Please,Login."})
		}
	}
}
