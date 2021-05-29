//+build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Runs go mod download and then installs the binary.
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	err := runDocker()
	if err != nil {
		return err
	}
	ch := make(chan error)
	go func() {
		ch <- sh.RunV("go", "run", "../cmd/zone-master/zone-master.go", "--config", "../configs/zone-master.yml", "serve")
	}()

	go func() {
		ch <- sh.RunV("go", "run", "../cmd/world-master/world-master.go", "--config", "../configs/world-master.yml", "serve")
	}()

	go func() {
		ch <- sh.RunV("go", "run", "../cmd/world/world.go", "--config", "../configs/world.yml", "serve")
	}()

	go func() {
		ch <- sh.RunV("go", "run", "../cmd/login/login.go", "--config", "../configs/login.yml", "serve")
	}()

	go func() {
		ch <- sh.RunV("go", "run", "../cmd/zone/zone.go", "--config", "../configs/zone.yml", "serve")
	}()

	<-ch
	return nil
}

func runDocker() error {
	return sh.Run("docker-compose", "-f", "../docker-compose.dev.yml", "up", "-d")
}
