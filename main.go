package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wingsxdu/tinyurl/server"
)

func main() {
	err := os.Mkdir("./log", os.ModePerm)
	f, err := os.Create("./log/httpsWarn.log")
	if err != nil {
		panic(err)
	}
	server.InitServer()
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

	e.GET("/t/:tinyUrl", getUrl)
	e.POST("/t", postUrl)
	e.PUT("/t", putUrl)
	e.DELETE("/t", deleteUrl)

	fmt.Printf("PID is：%d", os.Getpid())
	e.Logger.Warn(e.Start(":80"))
}

// http://localhost/t/2n9d
func getUrl(c echo.Context) error {
	tinyUrl := c.Param("tinyUrl")
	url, err := server.GetTinyUrl(tinyUrl)
	if err != nil {
		c.Error(err)
	}
	// if the tinyUrl doesn't exist
	if url == nil {
		return c.String(http.StatusAccepted, "The tinyUrl doesn't exist")
	}
	return c.Redirect(http.StatusTemporaryRedirect, string(url))
}

// http://localhost/t?url=https://www.google.com/
func postUrl(c echo.Context) error {
	url := c.QueryParam("url")
	index, err := server.PostTinyUrl([]byte(url))
	if err != nil {
		c.Error(err)
	}
	fmt.Println(c.Request().Host)
	tinyUrl := c.Request().Host + "/t/" + index
	return c.String(http.StatusOK, tinyUrl)
}

// http://localhost/t?tinyurl=2n9d&newurl=https://cn.bing.com/
func putUrl(c echo.Context) error {
	tinyUrl := c.QueryParam("tinyurl")
	newUrl := c.QueryParam("newurl")
	err := server.PutTinyUrl(tinyUrl, newUrl)
	if err != nil {
		c.Error(err)
	}
	return c.String(http.StatusOK, "Update tinyUrl Successfully")
}

// http://localhost/t?tinyurl=2n9d
func deleteUrl(c echo.Context) error {
	tinyUrl := c.QueryParam("tinyurl")
	err := server.DeleteTinyUrl(tinyUrl)
	if err != nil {
		c.Error(err)
	}
	return c.String(http.StatusOK, "Delete tinyUrl Successfully")
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
