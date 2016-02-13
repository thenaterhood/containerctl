package systemd

import (
  "os/exec"
  "fmt"
)

func RunMachinectlCmd(cmd, container string) {
    _, err := exec.Command("machinectl", cmd, container).Output()

    if err != nil {
        fmt.Println(cmd + " failed on " + container)
        fmt.Println(err)
    }
}
