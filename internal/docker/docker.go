package docker

import "os/exec"

func Available() bool {
	return exec.Command("docker", "version").Run() == nil
}
