package api

import "time"

func createErrorMsg(err error) map[string]string {
	return map[string]string{
		"error":     err.Error(),
		"timestamp": time.Now().Format(time.RFC3339),
	}
}
