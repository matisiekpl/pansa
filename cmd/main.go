package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/aymerick/raymond"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matisiekpl/pansa-plan/internal/controller"
	"github.com/matisiekpl/pansa-plan/internal/dto"
	"github.com/matisiekpl/pansa-plan/internal/repository"
	"github.com/matisiekpl/pansa-plan/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

func main() {
	raymond.RegisterHelper("inc", func(val1 int) string {
		return strconv.Itoa(val1 + 1)
	})

	if err := godotenv.Load(); err != nil {
		logrus.Info("no .env file loaded")
	}

	config := dto.NewConfig()
	repositories := repository.NewRepositories(config)
	services := service.NewServices(repositories)
	controllers := controller.NewControllers(services)

	e := echo.New()
	e.Use(middleware.Logger())
	e.HidePort = true
	e.HideBanner = true

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"service": "pansa-api", "contact": "Mateusz Woźniak <mateusz@mateuszwozniak.com>"})
	})
	e.GET("/report", controllers.Report().Generate)
	e.GET("/aip", controllers.Publication().Index)
	e.GET("/docs", controllers.Scalar().UI)
	e.GET("/docs/openapi.yaml", controllers.Scalar().Spec)
	if config.NotamEnabled() {
		e.GET("/notam/:icao", controllers.Notam().Index)
	} else {
		logrus.Warn("notam feature disabled: missing FAA_KEY/FAA_SECRET")
		e.GET("/notam/:icao", func(c echo.Context) error {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "notam feature disabled"})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		logrus.Infof("starting lambda gateway")
		lambda.Start(echoadapter.New(e).ProxyWithContext)
	} else {
		logrus.Infof("listening on port: %s", port)
		logrus.Fatal(e.Start(":" + port))
	}
}
