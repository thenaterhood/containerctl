package containerops

import(
  "fmt"
  "path"
  "os/exec"
)

type DebianContainer struct {

  GenericContainer

}

func (c DebianContainer) aptInstall(pkg string) error {
    err := c.Exec("apt-get", "install", "-y", "--allow-unauthenticated", pkg)
    return err
}

func (c DebianContainer) Create() error {

    if c.Installed() {
      return fmt.Errorf("%s %s", c.Name(), "is already installed")
    }

    err := c.GenericContainer.Create()
    if err != nil {
      return err
    }

    dir := path.Join(c.Location(), c.Name())
    _, err = exec.Command("debootstrap", "sid", dir).Output()

    if err != nil {
      return err
    }

    c.createMachineId()

    err = c.aptInstall("dbus")
    if err != nil {
      return err
    }

    err = c.aptInstall("coreutils")
    if err != nil {
      return err
    }

    err = c.Exec("/lib/systemd/systemd-sysv-install", "enable", "dbus")
    return err
}
