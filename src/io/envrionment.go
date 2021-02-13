package io

import (
	"fmt"
	"os/exec"
)

//CheckRequiredPrograms check wether or not the programs are installed on the current machine
//Used by the script too ensure that the local environment is correctly setup
func CheckRequiredPrograms(programs []string) error {
	var errorHasOccuredFlag bool
	var programError error

	for _, program := range programs {
		if err := checkProgram(program); err != nil {
			errorHasOccuredFlag = true
			programError = err
		}
	}

	if errorHasOccuredFlag {
		return programError
	}

	return nil
}

func checkProgram(program string) error {
	cmdProgramVersion := exec.Command(program, "--version")

	programVersion, err := cmdProgramVersion.Output()
	if err != nil {
		PrintRedCross(fmt.Sprintf("%s is not installed or could not be found in path\r\n", program))

		return err
	}

	PrintGreenCheckmark(fmt.Sprintf("%s is installed with version: %s", program, programVersion))

	return nil
}
