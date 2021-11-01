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
		ch <- sh.RunV("go", "run", "../cmd/zone-master/zone-master.go", "--config", "../configs/local/zone-master.yml", "serve")
	}()

	go func() {
		ch <- sh.RunV("go", "run", "../cmd/world-master/world-master.go", "--config", "../configs/local/world-master.yml", "serve")
	}()

	go func() {
		ch <- sh.RunV("go", "run", "../cmd/world/world.go", "--config", "../configs/local/world.yml", "serve")
	}()

	go func() {
		ch <- sh.RunV("go", "run", "../cmd/login/login.go", "--config", "../configs/local/login.yml", "serve")
	}()

	go func() {
		ch <- sh.RunV("go", "run", "../cmd/zone/zone.go", "--config", "../configs/local/zone.yml", "serve")
	}()

	<-ch
	return nil
}

func runDocker() error {
	return sh.Run("docker-compose", "-f", "../configs/local/docker-compose.yml", "up", "-d")
}
