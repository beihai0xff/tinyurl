package echoapi

import (
	"github.com/labstack/echo/v4"
	"github.com/wingsxdu/tinyurl"
	"net/http"
)

func New() {
	tinyurl.New()
}

// GET http://localhost/t/2n9d
func GetUrl(c echo.Context) error {
	tinyUrl := c.Param("tinyUrl")
	url, err := tinyurl.Get(tinyUrl)
	if err != nil {
		c.Error(err)
	} else if url == nil { // if the tinyUrl doesn't exist
		return c.String(http.StatusAccepted, "The tinyUrl doesn't exist")
	}
	return c.Redirect(http.StatusTemporaryRedirect, string(url))
}

// GET http://localhost/gett/2n9d
func Gett(c echo.Context) error {
	tinyUrl := c.Param("tinyUrl")
	url, err := tinyurl.Get(tinyUrl)
	if err != nil {
		c.Error(err)
	} else if url == nil { // if the tinyUrl doesn't exist
		return c.String(http.StatusAccepted, "The tinyUrl doesn't exist")
	}
	return c.String(http.StatusOK, string(url))
}

// POST http://localhost/t?url=https://www.google.com/
func PostUrl(c echo.Context) error {
	url := c.QueryParam("url")
	index, err := tinyurl.Create([]byte(url))
	if err != nil {
		c.Error(err)
	}
	tinyUrl := c.Request().Host + "/t/" + index
	return c.String(http.StatusOK, tinyUrl)
}

// PUT http://localhost/t?tinyurl=2n9d&newurl=https://cn.bing.com/
func PutUrl(c echo.Context) error {
	tinyUrl := c.QueryParam("tinyurl")
	newUrl := c.QueryParam("newurl")
	err := tinyurl.Update(tinyUrl, newUrl)
	if err != nil {
		c.Error(err)
	}
	return c.String(http.StatusOK, "Update tinyUrl Successfully")
}

// DELETE http://localhost/t?tinyurl=2n9d
func DeleteUrl(c echo.Context) error {
	tinyUrl := c.QueryParam("tinyurl")
	err := tinyurl.Delete(tinyUrl)
	if err != nil {
		c.Error(err)
	}
	return c.String(http.StatusOK, "Delete tinyUrl Successfully")
}
