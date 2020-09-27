//+build mage

package main

import (
	"fmt"
	"os"
	"github.com/magefile/mage/sh"
)


// Runs go mod download and then installs the binary.
func Build() error {
	cmds := []string{
		"login",
		"world",
		"world-master",
		"zone",
		"zone-master",
	}

	err := runDocker()
	if err != nil {
		return err
	}

	for _, c := range cmds {
		err := buildCmd(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func runDocker() error {
	return sh.Run("docker-compose","up", "-d")
}

//func vendor() error {
//	return sh.RunV("go", "mod", "vendor", "../")
//}

func buildCmd(name string) error {
	var out string
	workdir := fmt.Sprintf("../cmd/%v/", name)
	if isWindows() {
		out = fmt.Sprintf("./package/%v/%v", name, name+".exe")
	} else {
		out = fmt.Sprintf("./package/%v/%v", name, name+".exe")
	}

	// TODO: build for linux
	return sh.RunV("go", "build", "-mod", "mod",  "--race", "-o", out, workdir)
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
