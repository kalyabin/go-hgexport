package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

type File struct {
	Path   string
	status string
}

// check if file exists
func (file *File) IsExists() bool {
	if file.Path == "" {
		return false
	}

	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		return false
	}

	return true
}

// returns false if file can't be added to archive
func (file *File) CanAdd() bool {
	return file.IsExists() && file.status != "D" && file.status != "R"
}

// add file to zip archive
func (file *File) AddToArchive(zipWriter *zip.Writer) error {
	fileReader, err := os.Open(file.Path)

	if err != nil {
		return fmt.Errorf("Can't open file %s: %s", file.Path, err)
	}

	defer fileReader.Close()

	fileWriter, err := zipWriter.Create(file.Path)

	if err != nil {
		return fmt.Errorf("Can't add file %s to archive: %s", file.Path, err)
	}

	if _, err := io.Copy(fileWriter, fileReader); err != nil {
		return fmt.Errorf("Can't writer file %s to archive: %s", file.Path, err)
	}

	return nil
}
