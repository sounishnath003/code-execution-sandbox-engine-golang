package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sounishnath/code-sandbox-runner/cmd/api"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Gzip())

	api.InitializeContainerPool(25)

	fmt.Println("server is up and running on http://localhost:3000")
	e.POST("/api/submit", api.ExecuteCodeHandler)

	if err := e.Start(":3000"); err != nil {
		e.Logger.Info("shutting down the server")
	}

}
