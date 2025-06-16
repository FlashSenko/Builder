package build

import (
	"fmt"
	"os/exec"
	"strings"
)

func execute(command string, sourceArg string, destArg string) error {
	command = strings.Replace(strings.Replace(command, "<s", sourceArg, 1), "<S", sourceArg, 1)
	command = strings.Replace(strings.Replace(command, "<d", destArg, 1), "<D", destArg, 1)

	output, err := exec.Command("cmd", "/C", command).CombinedOutput()
	if err != nil {
		fmt.Printf("%s", output)
		return err
	}

	return nil
}
