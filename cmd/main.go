package main

import (
	"fmt"
	"sync"

	"github.com/labstack/echo"
	"github.com/sounishnath/code-sandbox-runner/cmd/api"
)

func main() {
	e := echo.New()

	api.InitializeContainerPool(5)
	api.JobQueue = make(chan api.Job, 10)

	// Start the worker pool.
	numWorkers := 5
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go api.Worker(i, &wg)
	}

	fmt.Println("server is up and running on http://localhost:3000")
	e.POST("/api/submit", api.ExecuteCodeHandler)

	// Start the server (in a separate goroutine if you want to wait on workers).
	go func() {
		if err := e.Start(":3000"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	wg.Wait()
}
