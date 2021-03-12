package embed

import (
	"os"
	"path/filepath"
	"strings"
)

func filePaths(root string) (map[string]string, error) {
	collection := make(map[string]string)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path == root {
			return nil
		} else if info.IsDir() {
			subCollection, err := filePaths(path)
			if err != nil {
				return err
			}
			for sc, _ := range subCollection {
				collection[sc] = ""
			}
			return nil
		} else if strings.HasSuffix(info.Name(), ".gohtml") == false {
			return nil
		}

		collection[path] = ""

		return nil
	})

	return collection, err
}
