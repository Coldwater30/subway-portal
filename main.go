package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	wdGroup := e.Group("/wifidog")
	wdGroup.GET("/ping/", ping)
	wdGroup.GET("/login/", login)
	wdGroup.GET("/logincheck", loginCheck)
	wdGroup.GET("/auth/", auth)
	wdGroup.GET("/portal/", portal)
	e.Logger.Fatal(e.Start(":8082"))
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "Pong!")
}

func login(c echo.Context) error {
	return 
}