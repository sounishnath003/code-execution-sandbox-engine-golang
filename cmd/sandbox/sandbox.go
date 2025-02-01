package sandbox

import "errors"

// RunCodeInDocker executes the provided code into the docker container.
//
// Returns StdErr, StdOut, Error.
func RunCodeInDocker(language, base64EncodedCodeString string) (string, string, error) {
	switch language {
	case "python3":
		return pythonCodeRunner(base64EncodedCodeString)
	case "java":
		return javaCodeRunner(base64EncodedCodeString)

	default:
		return "", "", errors.New("Language not supported")
	}
}
