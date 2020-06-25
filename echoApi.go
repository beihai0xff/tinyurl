package tinyurl

import (
	"github.com/labstack/echo/v4"
	"github.com/wingsxdu/tinyurl/server"
	"net/http"
)

func New() {
	server.InitServer()
}

// http://localhost/t/2n9d
func GetUrl(c echo.Context) error {
	tinyUrl := c.Param("tinyUrl")
	url, err := server.GetTinyUrl(tinyUrl)
	if err != nil {
		c.Error(err)
	} else if url == nil { // if the tinyUrl doesn't exist
		return c.String(http.StatusAccepted, "The tinyUrl doesn't exist")
	}
	return c.Redirect(http.StatusTemporaryRedirect, string(url))
}

// http://localhost/t?url=https://www.google.com/
func PostUrl(c echo.Context) error {
	url := c.QueryParam("url")
	index, err := server.PostTinyUrl([]byte(url))
	if err != nil {
		c.Error(err)
	}
	tinyUrl := c.Request().Host + "/t/" + index
	return c.String(http.StatusOK, tinyUrl)
}

// http://localhost/t?tinyurl=2n9d&newurl=https://cn.bing.com/
func PutUrl(c echo.Context) error {
	tinyUrl := c.QueryParam("tinyurl")
	newUrl := c.QueryParam("newurl")
	err := server.PutTinyUrl(tinyUrl, newUrl)
	if err != nil {
		c.Error(err)
	}
	return c.String(http.StatusOK, "Update tinyUrl Successfully")
}

// http://localhost/t?tinyurl=2n9d
func DeleteUrl(c echo.Context) error {
	tinyUrl := c.QueryParam("tinyurl")
	err := server.DeleteTinyUrl(tinyUrl)
	if err != nil {
		c.Error(err)
	}
	return c.String(http.StatusOK, "Delete tinyUrl Successfully")
}
