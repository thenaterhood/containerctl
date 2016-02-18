package containerops

import(
  "fmt"
  "path"
  "os/exec"
)

type ArchContainer struct {

  GenericContainer

}

// Installs a base Arch install into a container directory using pacstrap
func (c ArchContainer) Create() error {

    if c.Installed() {
        return fmt.Errorf("%s %s", c.Name(), "is already installed")
    }

    err := c.GenericContainer.Create()

    if err != nil {
      return err
    }

    dir := path.Join(c.Location(), c.Name())
    _, err = exec.Command("pacstrap", "-c", "-d", dir, "base", "--ignore", "linux").Output()

    return err
}
