package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

func pythonCodeRunner(base64EncodedCodeString string) (string, string, error) {
	// Prepare the container runtime.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// cmdString := fmt.Sprintf(`echo "%s" | base64 -d | go run /dev/stdin`, base64EncodedCodeString)
	cmdString := fmt.Sprintf(`echo "%s" | base64 -d | python3`, base64EncodedCodeString)

	// Run the code inside docker container
	cmd := exec.CommandContext(ctx, "docker", "run", "python:3.11-slim", "sh", "-c", cmdString)

	// Store the cmd output/error into buffers.
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), out.String(), err
	}

	return stderr.String(), out.String(), nil

}
