package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Repo struct {
	CmdPath string
	RevFrom string
	RevTo   string
}

// run some hg command with some params
func (instance *Repo) run(params ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(instance.CmdPath, params...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	commandTxt := instance.CmdPath + " " + strings.Join(params, " ")

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Failed to run command %s: %s", commandTxt, err)
	}

	return string(stdout.Bytes()), nil
}

// check if repository exists at current path
func (instance *Repo) CheckRepo() error {
	if _, err := instance.run("root"); err != nil {
		return err
	}

	return nil
}

// get all changed files between revisions
func (instance *Repo) GetChangedFiles() ([]File, error) {
	var result []File

	filesList, err := instance.run("status", "--rev", instance.RevFrom+":"+instance.RevTo)

	if err != nil {
		return result, err
	}

	for _, line := range strings.Split(filesList, "\n") {
		if line != "" {
			pieces := strings.SplitN(line, " ", 2)
			if len(pieces) == 2 {
				file := File{Path: pieces[1], status: pieces[0]}
				result = append(result, file)
			}
		}
	}

	return result, nil
}

