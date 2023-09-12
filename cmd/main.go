package main

import (
	"github.com/aymerick/raymond"
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/pansa-plan/internal/controller"
	"github.com/matisiekpl/pansa-plan/internal/repository"
	"github.com/matisiekpl/pansa-plan/internal/service"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {
	raymond.RegisterHelper("inc", func(val1 int) string {
		return strconv.Itoa(val1 + 1)
	})

	repositories := repository.NewRepositories()
	services := service.NewServices(repositories)
	controllers := controller.NewControllers(services)

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	e.GET("/report", controllers.Report().Generate)
	e.GET("/aip", controllers.Publication().Index)
	e.GET("/notam/:icao", controllers.Notam().Index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("listening on port: %s", port)
	logrus.Fatal(e.Start(":" + port))
}
