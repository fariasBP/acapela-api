package controllers

import "github.com/labstack/echo/v4"

type dat struct {
	Appname string `json:"appname"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
}

func InfoWeb(c echo.Context) error {

	u := &dat{
		Appname: "Acapela",
		Code:    200,
		Msg:     "Hello World!!!",
	}

	return c.JSON(200, u)
}
