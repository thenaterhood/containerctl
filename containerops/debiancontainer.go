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

// Specialized install function for Debian containers.
// Runs debootstrap for a base install, then manually installs
// dbus, coreutils, and a machine-id into the container to
// make it usable as a standalone container.
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

    err = c.createMachineId()
    if err != nil {
      return err
    }

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
