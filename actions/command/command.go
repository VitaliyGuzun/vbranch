package command

import (
	"os"
	"os/exec"
)

func Run(arg ...string) error {
	command := exec.Command(arg[0], arg[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()

	return err
}
