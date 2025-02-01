package sandbox

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// RunCodeInDocker executes the provided code into the docker container.
func RunCodeInDocker(language, code string) {
	containerName := "code-runner"

	// Prepare the container runtime.
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	cmdString := fmt.Sprintf(`docker run --rm --name %s golang:1.23-alpine sh -c echo '%s' | go run /dev/stdin`, containerName, code)
	cmd := exec.CommandContext(ctx, cmdString)
}