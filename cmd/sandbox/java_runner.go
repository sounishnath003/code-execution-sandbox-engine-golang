package sandbox

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"os"
	"os/exec"
)

// javaCodeRunner used to run the Java code into the sandbox environment
func javaCodeRunner(base64EncodedCode string) (string, string, error) {
	containerName := fmt.Sprintf("java-sandbox-%d", time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Decode the base64-encoded code
	code, err := base64.StdEncoding.DecodeString(base64EncodedCode)
	if err != nil {
		return err.Error(), "", err
	}

	tempFile, err := os.Create(fmt.Sprintf("%d.java", time.Now().UnixNano()))
	if err != nil {
		return "", "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	// Write the code to the temporary file
	if _, err := tempFile.Write([]byte(code)); err != nil {
		return "", "", fmt.Errorf("failed to write code to temporary file: %w", err)
	}

	// Build the Docker image with the compiled class
	dockerfileName := "Dockerfile"
	dockerfileContent := fmt.Sprintf(`
		FROM openjdk:11-jdk-slim
		WORKDIR /app
		COPY %s HelloWorld.java
		ENTRYPOINT ["java", "HelloWorld.java"]
		`, tempFile.Name())

	if err := os.WriteFile(dockerfileName, []byte(dockerfileContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write Dockerfile: %w", err)
	}
	defer os.Remove(dockerfileName)  // Cleanup Dockerfile
	defer os.Remove(tempFile.Name()) // Cleanup Dockerfile

	// Build the image from the Dockerfile
	cmd := exec.CommandContext(ctx, "docker", "build", "-t", containerName, ".")
	var buildOut, buildErr bytes.Buffer
	cmd.Stdout = &buildOut
	cmd.Stderr = &buildErr
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Dockerfile build error: %s", buildErr.String()), "", err
	}

	// Run the container to execute the compiled class
	cmd = exec.CommandContext(ctx, "docker", "run", "--rm", containerName)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stderr.String(), out.String(), err
	}

	// Run the container to execute the compiled class
	cmd = exec.CommandContext(ctx, "docker", "rmi", "-f", containerName)
	if err := cmd.Run(); err != nil {
		return "", "", err
	}

	return stderr.String(), out.String(), nil
}
