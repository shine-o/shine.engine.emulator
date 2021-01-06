package main

import (
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"strconv"
	"strings"
)

type NcDepartment uint16

const (
	NC_DEPARTMENT_COMMAND_1 = 2500
	NC_DEPARTMENT_COMMAND_2 = 2501
)

func main() {

	// load network-commands.yml file

	// for each command in department
	// 		const (
	//	  	  NC_TYPE_CONCRETE NcType =
	//	 	)

	// create the opcode from department + command

	//
	pcl, err := networking.LoadCommandList("configs/network-commands.yml")

	if err != nil {
		logger.Error(err)
	}

	res := "const ( %v \n)\n"
	commands := ""
	for _, v1 := range pcl.Departments {
		dept, err := strconv.ParseInt(stripRawHex(v1.HexID), 16, 64)
		if err != nil {
			logger.Error(err)
		}

		for k2, v2 := range v1.ProcessedCommands {
			command, err := strconv.ParseInt(stripRawHex(k2), 16, 64)
			if err != nil {
				logger.Error(err)
			}
			//   command = (department << 10) | command & 0x3FF;
			opcode := dept<<10 | command&1023
			commands += fmt.Sprintf("%v = %v\n", v2, opcode)
		}
	}
	fmt.Printf(res, commands)
}

func stripRawHex(rh string) string {
	rh = strings.ReplaceAll(rh, "\n", "")
	rh = strings.ReplaceAll(rh, " ", "")
	rh = strings.ReplaceAll(rh, "0x", "")
	rh = strings.ReplaceAll(rh, "\t", "")
	return rh
}
