package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"kos-api/domains/node"
	"kos-api/domains/space"
	"kos-api/internal/config"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	nodeModule, err := node.New(appConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer nodeModule.Close()

	spaceModule, err := space.New(appConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer spaceModule.Close()

	e := echo.New()
	e.HideBanner = true
	e.GET("/api", healthHandler)
	nodeModule.RegisterRoutes(e)
	spaceModule.RegisterRoutes(e)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "KOS API root")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(e.Start(":" + port))
}

type healthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, healthResponse{
		Status:    "ok",
		Service:   "kos-api",
		Timestamp: time.Now().UTC(),
	})
}
