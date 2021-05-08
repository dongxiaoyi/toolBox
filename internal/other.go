package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"errors"
)

func AbsPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		panic(errors.New(`error: Can't find "/" or "\".`))
		return ""
	}
	return path[0 : i+1]
}