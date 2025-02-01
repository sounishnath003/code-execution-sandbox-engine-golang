package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func ExecuteCodeHandler(c echo.Context) error {
	var req CodeSubmission
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Create a result channel for this job.
	jobResultChan := make(chan JobResult, 1)

	// Create a new job.
	job := Job{
		ID:       generateJobID(), // You can implement this function to generate unique IDs.
		Language: req.Language,
		Code:     req.Base64EncodedCode,
		Result:   jobResultChan,
	}

	// Submit the job into the job queue.
	JobQueue <- job

	// Wait for the job to finish.
	result := <-jobResultChan

	execResult := ExecutionResult{JobResult: result}
	if result.Err != nil {
		if result.StdErr != "" {
			execResult.JobResult.StdErr = result.StdErr
		} else {
			execResult.JobResult.StdErr = result.Err.Error()
		}
	}

	return c.JSON(http.StatusOK, execResult)
}

// generateJobID is a simple helper to generate a job identifier.
func generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}
