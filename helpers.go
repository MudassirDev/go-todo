package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

func validateContentType(header http.Header, expectedType string) error {
	if header.Get("Content-Type") == expectedType {
		return nil
	}
	return fmt.Errorf("invalid content type")
}

func parseTemplates() {
	filepath.WalkDir("static/templates", func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".html") || d.IsDir() {
			return nil
		}
		Templates.ParseFiles(path)
		return nil
	})
}
