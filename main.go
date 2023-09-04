package main

import (
	"github.com/aymerick/raymond"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

func main() {
	raymond.RegisterHelper("inc", func(val1 int) string {
		return strconv.Itoa(val1 + 1)
	})
	e := echo.New()
	e.HideBanner = true
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.POST("/plan", generatePlan)
	e.GET("/report", generateReport)
	e.GET("/aip", serveAIP)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("listening on port: %s", port)
	logrus.Fatal(e.Start(":" + port))
}
