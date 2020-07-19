package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

const (
	checkCommand = "/opt/resource/check"
)

type CheckCommand struct {
	Run CheckRunCommand `command:"run" description:"Runs the '/opt/resource/check' command on a resource"`
}

type CheckRunCommand struct {
	Input  string `long:"input" short:"i" required:"true" description:"The path to a JSON file which will be streamed to STDIN of the '/opt/resource/check' executable"`
	Output string `long:"output" short:"o" required:"true" description:"The path to store the STDOUT of the '/opt/resource/check' executable"`
}

func (crc *CheckRunCommand) Execute([]string) error {
	checkInput, err := ioutil.ReadFile(crc.Input)
	if err != nil {
		return errors.New("reading resource check input file")
	}

	command := exec.Command(checkCommand)

	stdin, err := command.StdinPipe()
	if err != nil {
		return fmt.Errorf("obtaining stdin of the %v command", checkCommand)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("obtaining stdout of the %v command", checkCommand)
	}

	_, err = stdin.Write(checkInput)
	if err != nil {
		return fmt.Errorf("writing to the stdin of the %v command", checkCommand)
	}

	err = stdin.Close()
	if err != nil {
		return fmt.Errorf("closing the stdin of the %v command", checkCommand)
	}

	err = command.Start()
	if err != nil {
		return fmt.Errorf("starting the %v command", checkCommand)
	}

	checkOutput, err := ioutil.ReadAll(stdout)
	if err != nil {
		return fmt.Errorf("reading the stdout of the %v command", checkCommand)
	}

	err = ioutil.WriteFile(crc.Output, checkOutput, 0644)
	if err != nil {
		return errors.New("writing resource check output file")
	}

	command.Wait()

	return nil
}
