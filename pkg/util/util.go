package util

import (
	"io"
	"os"
	"path/filepath"
)

func SaveContentToFile(reader io.Reader, fileName string) (io.Reader, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(filepath.Join(dir, "html"), 0755)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filepath.Join(dir, "html", fileName+".html"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.TeeReader(reader, file), nil
}
