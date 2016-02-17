package systemd

import (
  "os/exec"
)

func RunMachinectlCmd(cmd, container string) error {
    _, err := exec.Command("machinectl", cmd, container).Output()

    return err
}
