package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Gateway a struct type that records gateway address and port
type Gateway struct {
	GWAddress string
	GWPort    string
}

func main() {
	// create echo instances and register template renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("web/*.html")),
	}
	e := echo.New()
	e.Renderer = t
	// registers a new route with path prefix to serve static files
	// from the web directory
	e.Static("/web", "web")
	// configure logging format for each request
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\n",
	}))
	// register handlers for url route
	wdGroup := e.Group("/wifidog") //group routes with prefix /wifidog
	wdGroup.GET("/ping/", Ping)
	wdGroup.GET("/login/", Login)
	wdGroup.POST("/logincheck", LoginCheck)
	wdGroup.GET("/auth/", Auth)
	wdGroup.GET("/portal/", Portal)

	e.Logger.Fatal(e.Start(":8082"))
}

func generateToken(mac string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(mac + strconv.FormatFloat(rand.Float64(), 'e', 6, 32)))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
func index(c echo.Context) error {
	return c.Render(http.StatusOK, "welcome.html", "World")
}

// Login protocol module
func Login(c echo.Context) error {
	gwAddress := c.QueryParam("gw_address")
	gwPort := c.QueryParam("gw_port")
	gateway := Gateway{
		GWAddress: gwAddress,
		GWPort:    gwPort,
	}
	return c.Render(http.StatusOK, "login", gateway)
}

// LoginCheck check if the password is right, generate token by email and redirect
// user to gateway with token
func LoginCheck(c echo.Context) error {
	password := c.FormValue("password")
	if password != "12345678" {
		log.Println("wrong password")
		return c.String(http.StatusOK, "wrong password!")
	}

	gwAddress := c.FormValue("gw_address")
	gwPort := c.FormValue("gw_port")
	email := c.FormValue("email")
	token := generateToken(email)
	uri := fmt.Sprintf("http://%s:%s/wifidog/auth?token=%s", gwAddress, gwPort, token)
	return c.Redirect(http.StatusPermanentRedirect, uri)
}

// Auth protocol module for token verification and user status bookeeping
func Auth(c echo.Context) error {
	stage := c.QueryParam("stage")

	switch stage {
	case "login":
		return c.String(http.StatusOK, "Auth: 1")
	case "counters":
		return c.String(http.StatusOK, "Auth: 1")
	default:
		return c.String(http.StatusOK, "OK")
	}
}

// Ping protocol module, for gateway to check if server alive or not
func Ping(c echo.Context) error {
	ans := "Pong!"
	log.Println("Answer: ", ans)
	return c.String(http.StatusOK, ans)
}

// Portal This is to return the page after successful login
func Portal(c echo.Context) error {
	return c.Render(http.StatusOK, "welcome.html", nil)
}
