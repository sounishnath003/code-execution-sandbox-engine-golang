package api

type CodeSubmission struct {
	Language          string `json:"language"`
	Base64EncodedCode string `json:"base64EncodedCode"`
}

type ExecutionResult struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr,omitempty"`
}
