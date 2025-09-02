package main

import (
	"fmt"
	"net/http"
)

func validateContentType(header http.Header, expectedType string) error {
	if header.Get("Content-Type") == expectedType {
		return nil
	}
	return fmt.Errorf("invalid content type")
}
