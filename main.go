package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
)

func packFilesToArchive(filesList []File, archiveName string) {
	zipFile, err := os.OpenFile(archiveName, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("Can't open zip file: %s", err)
	}

	defer zipFile.Close()

	buffer := new(bytes.Buffer)

	writer := zip.NewWriter(buffer)

	for _, file := range filesList {
		if !file.CanAdd() {
			fmt.Printf("File %s can't be added to archive\n", file.Path)
			continue
		}

		if err := file.AddToArchive(writer); err != nil {
			log.Fatalf("Error: %s", err)
		}

		fmt.Println(file.Path)
	}
}

// get revision from, revision to and archive file name
func getArgs() (string, string, string, error) {
	var revFrom, revTo, archiveName string

	flag.StringVar(&revFrom, "from", "", "Repo revision to start files list")
	flag.StringVar(&revTo, "to", "", "Repo revision to end files list")
	flag.StringVar(&archiveName, "out", "", "Result ZIP archive name")

	flag.Usage = func() {
		fmt.Println("Hgexport - the command line tool to export all changed files between some Mercurial repo revisions")
		flag.PrintDefaults()
	}

	flag.Parse()

	if revFrom == "" || revTo == "" || archiveName == "" {
		flag.Usage()
		return revFrom, revTo, archiveName, fmt.Errorf("Empty arguments\n")
	}

	return revFrom, revTo, archiveName, nil
}

func main() {
	revFrom, revTo, archiveName, err := getArgs()
	if err != nil {
		return
	}

	hg := Repo{CmdPath: "hg", RevFrom: revFrom, RevTo: revTo}

	if err := hg.CheckRepo(); err != nil {
		log.Fatalf("Repo not found here: %s", err)
	}

	changedFiles, err := hg.GetChangedFiles()

	if err != nil {
		log.Fatalf("Can't get changed files list: %s", err)
	}

	if len(changedFiles) == 0 {
		fmt.Println("Empty files list")
		return
	}

	packFilesToArchive(changedFiles, archiveName)
}
