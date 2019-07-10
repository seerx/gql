package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

var commands = map[string]string{
	"windows": "cmd /c start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

// OpenURL calls the OS default program for uri
func OpenURL(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	cmd := exec.Command(run, uri)
	return cmd.Start()
}
