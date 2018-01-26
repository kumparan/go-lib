package sqlimporter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func getFileList(dir string) ([]string, error) {
	var files []string
	if err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f != nil && !f.IsDir() {
			// make sure that file have .sql extension, else ignore
			if strings.Contains(f.Name(), ".sql") {
				files = append(files, path)
			}
		}
		return err
	}); err != nil {
		return nil, err
	}
	return files, nil
}

// this required to split the string
var delimiter = "--end"

// parseFile will parse the entire .sql file to queries
func parseFiles(filepath string) ([]string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	queries := strings.Split(string(content), delimiter)
	return sanitizeQueries(queries), nil
}

// sanitize queries
// for now only skip queries if empty string or enter
func sanitizeQueries(queries []string) []string {
	r := strings.NewReplacer("\n", "")

	var sanitized []string
	for key := range queries {
		if queries[key] == "" || queries[key] == "\n" {
			continue
		}
		q := r.Replace(queries[key])
		q = strings.TrimSpace(q)
		sanitized = append(sanitized, q)
	}
	return sanitized
}
