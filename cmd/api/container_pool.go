package api

import (
	"fmt"
	"sync"

	"github.com/sounishnath/code-sandbox-runner/cmd/sandbox"
)

// Job represents a code execution task.
type Job struct {
	ID       string
	Language string
	Code     string
	Result   chan JobResult // channel to receive the job result
}

// JobResult represents the outcome of a job execution.
type JobResult struct {
	JobID       string `json:"jobid"`
	ContainerID string `json:"containerID"`
	StdOut      string `json:"stdout"`
	StdErr      string `json:"stderr,omitempty"`
	Err         error  `json:"error,omitempty"`
}

// Global job queue (buffered channel)
var JobQueue chan Job

// Container represents a running sandbox container.
type Container struct {
	ID string // For example, the Docker container ID
}

// Global container pool (buffered channel)
var containerPool chan *Container

// Global wait group allows to block a specific code block to
// allow a set of goroutines to complete execution.
var wg sync.WaitGroup

// Initialize the container pool with n containers.
func InitializeContainerPool(n int) {
	containerPool = make(chan *Container, n)
	for i := 0; i < n; i++ {
		// Here you would start a container via Docker API or CLI.
		// For this demo, we simply simulate with a dummy container ID.
		container := &Container{ID: fmt.Sprintf("container-%d", i+1)}
		containerPool <- container
	}
	
	// Initialize the buffer job queue can be replaced by Redis / other Queues
	JobQueue = make(chan Job, 2*n)
	wg.Wait()

	// Start the worker pool.
	numWorkers := 24
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg)
		fmt.Printf("worker %d has been registered inside container pool\n", i)
	}
}

// worker picks up jobs from the jobQueue.
func worker(workerID int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range JobQueue {
		fmt.Printf("Worker %d: Received job %s\n", workerID, job.ID)

		// Get an available container from the container pool.
		container := <-containerPool
		fmt.Printf("Worker %d: Using container %s for job %s\n", workerID, container.ID, job.ID)

		// Execute the job in the container.
		containerID, stdErr, stdOut, err := runCodeInContainer(container, job.Language, job.Code)
		// Return the container to the pool.
		containerPool <- container

		// Package the result.
		result := JobResult{
			JobID:       job.ID,
			ContainerID: containerID,
			StdOut:      stdOut,
			StdErr:      stdErr,
			Err:         err,
		}

		// Send the result back to the job's result channel.
		job.Result <- result
	}
}

// runCodeInContainer simulates executing the code inside a container.
// Replace the contents with your actual Docker execution logic.
func runCodeInContainer(container *Container, language, code string) (string, string, string, error) {
	stdErr, stdOut, err := sandbox.RunCodeInDocker(language, code)

	// For demonstration, we return dummy output.
	return container.ID, stdErr, stdOut, err
}
