//+build mage

package main

import (
	"fmt"
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

func buildCmd(name string) error {
	workdir := fmt.Sprintf("./cmd/%v/", name)
	out := fmt.Sprintf("./build/package/%v/%v", name, name)

	// TODO: build for linux

	return sh.RunV("go", "build", "--race", "-o", out, workdir)
}
