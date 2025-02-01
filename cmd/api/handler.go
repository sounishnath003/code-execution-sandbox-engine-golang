package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sounishnath/code-sandbox-runner/cmd/sandbox"
)

func ExecuteCodeHandler(c echo.Context) error {
	var req CodeSubmission

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, createErrorMsg(err))
	}

	stderr, stdout, err := sandbox.RunCodeInDocker(req.Language, req.Base64EncodedCode)
	result := ExecutionResult{StdOut: stdout}
	if err != nil {
		if len(stderr) > 0 {
			result.StdErr = stderr
		} else {
			result.StdErr = err.Error()
		}
	}

	return c.JSON(http.StatusOK, result)
}
