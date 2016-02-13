package systemd

import (
  "os/exec"
  "fmt"
)


func StartService(container string) (bool, string) {
    out, err := exec.Command("systemctl", "start", container).Output()

    if err != nil {
        fmt.Println("systemctl start failed on " + container)
        fmt.Println(err)
        return false, string(out)
    }

    return true, string(out)
}
