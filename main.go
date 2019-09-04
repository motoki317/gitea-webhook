package main

import (
	"github.com/motoki317/gitea-webhook/webhook"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server successfully started!")
	})

	e.POST("/webhook", webhook.MakeWebhookHandler())

	err := e.Start(":8080")
	if err != nil {
		log.Println(err)
	}
}
