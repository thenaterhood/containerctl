package systemd

import (
  "os/exec"
  "fmt"
)


func StartService(container string) (error, string) {
    out, err := exec.Command("systemctl", "start", container).Output()

    if err != nil {
        err = fmt.Errorf("%s: %s %s", err, "systemctl start failed on", container)
        return err, string(out)
    }

    return nil, string(out)
}
