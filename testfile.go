package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/test", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "The application is up!")
	})
	e.Logger.Fatal(e.Start(":2000"))
}