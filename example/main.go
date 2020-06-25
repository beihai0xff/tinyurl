package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	. "github.com/wingsxdu/tinyurl"
)

func main() {
	err := os.Mkdir("./log", os.ModePerm)
	f, err := os.Create("./log/httpsWarn.log")
	if err != nil {
		panic(err)
	}
	New()
	/*	go func() {
		f2, err := os.Create("./log/httpWarn.log")
		if err != nil {
			panic(err)
		}
		h := echo.New()
		h.Use(middleware.Gzip())
		h.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: f2,
		}))
		h.GET("/t/:tinyUrl", getUrl)
		h.POST("/t", postUrl)
		h.PUT("/t/:tinyUrl", putUrl)
		h.DELETE("/t/:tinyUrl", deleteUrl)
		h.Logger.Warn(h.Start(":80"))
	}()*/

	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: f,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.GET("/t/:tinyUrl", GetUrl)
	e.POST("/t", PostUrl)
	e.PUT("/t", PutUrl)
	e.DELETE("/t", DeleteUrl)

	fmt.Printf("PID is：%d", os.Getpid())
	e.Logger.Warn(e.Start(":80"))
}

type httpError struct {
	code int
	Key  string `json:"error"`
	Msg  string `json:"message"`
}

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

	var res = httpError{code: http.StatusInternalServerError, Key: "InternalServerError"}

	if he, ok := err.(*echo.HTTPError); ok {
		res.code = he.Code
		res.Key = http.StatusText(res.code)
		res.Msg = err.Error()
	} else {
		res.Msg = http.StatusText(res.code)
	}

	if !c.Response().Committed {
		err := c.JSON(res.code, res)
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
