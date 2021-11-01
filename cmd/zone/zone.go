package main

import (
	"github.com/pkg/profile"
	"github.com/shine-o/shine.engine.emulator/pkg/zone/cmd"
)

func main() {
	defer profile.Start(profile.MemProfile).Stop()
	cmd.Execute()
}
