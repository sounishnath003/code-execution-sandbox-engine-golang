package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/sounishnath/code-sandbox-runner/cmd/api"
)

func main() {
	e := echo.New()

	fmt.Println("server is up and running on http://localhost:3000")
	e.POST("/api/submit", api.ExecuteCodeHandler)

	e.Start(":3000")
}
