package containerops

import(
  "fmt"
  "os"
  "path"
  "os/exec"
)

type DebianContainer struct {

  GenericContainer

}

func (c DebianContainer) New(name, location, uuid string, installed bool) {
  c.GenericContainer.name = name
  c.GenericContainer.location = location
  c.GenericContainer.uuid = uuid
  c.GenericContainer.installed = installed
}

func (c DebianContainer) aptInstall(pkg string) {
    c.Exec("apt-get", "install", "-y", "--allow-unauthenticated", pkg)
}

func (c DebianContainer) Create() {

    c.GenericContainer.Create()

    if c.Installed() {
        fmt.Println(c.Name() + " is already installed")
        os.Exit(1)
    }

    dir := path.Join(c.Location(), c.Name())
    _, err := exec.Command("debootstrap", "sid", dir).Output()

    c.createMachineId()
    c.aptInstall("dbus")
    c.aptInstall("coreutils")
    c.Exec("/lib/systemd/systemd-sysv-install", "enable", "dbus")

    if err != nil {
        fmt.Println(err)
    }
}
