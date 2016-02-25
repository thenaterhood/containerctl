package containerops

import(
    "fmt"
    "os/exec"
    "os"
    "path"
    "strings"
    "time"
    "github.com/thenaterhood/containerctl/systemd"
    "github.com/thenaterhood/containerctl/system"
)

// A generic container type, implementing basic methods for handling
// systemd-nspawn containers.
type GenericContainer struct {
    Container

    name string
    location string
    installed bool
    uuid string
}

func (c GenericContainer) Name() string {
  return c.name
}

func (c GenericContainer) Location() string {
  return c.location
}

func (c GenericContainer) Installed() bool {
  return c.installed
}

func (c GenericContainer) Uuid() string {
  return c.uuid
}

func (c GenericContainer) Create() error {
    err := os.Mkdir(path.Join(c.Location(), c.Name()), 0755)
    return err
}

func (c GenericContainer) Destroy() error {
    c.Stop()
    time.Sleep(100 * time.Millisecond)

    dir := path.Join(c.Location(), c.Name())
    err := os.RemoveAll(dir)
    return err
}

func (c GenericContainer) createMachineId() error {
    dir := path.Join(c.Location(), c.Name())

    uuidbytes, err := exec.Command("uuidgen").Output()

    if err == nil {
      uuid := string(uuidbytes[:37])
      var machineid *os.File
      machineid, err = os.OpenFile(path.Join(dir, "etc", "machine-id"), os.O_WRONLY, 0600)
      defer machineid.Close()
      machineid.WriteString(strings.Replace(uuid, "-", "", -1))
    }

    return err
}

func (c GenericContainer) Exec(args ...string) error {

    dir := path.Join(c.Location(), c.Name())
    command := append([]string{"-D", dir}, args...)

    _, err := exec.Command("systemd-nspawn", command...).Output()
    if err != nil {
      err = fmt.Errorf("%s: %s %s", string(err.Error()), "systemd-nspawn", strings.Join(args, " "))
    }

    return err
}

func (c GenericContainer) Start() error {

  if ! c.Installed() {
    return fmt.Errorf("%s %s", c.Name(), "does not have a system installed.")
  }

  err := systemd.RunMachinectlCmd("start", c.Name())
  return err
}

func (c GenericContainer) Stop() error {
  if ! c.Installed() {
    return fmt.Errorf("%s %s", c.Name(), "does not have a system installed.")
  }

  err := systemd.RunMachinectlCmd("poweroff", c.Name())
  return err
}

func (c GenericContainer) UpdateUser(user *system.OSUser) error {
  err := user.UpdateEntry(path.Join(c.Location(), c.Name()))
  return err
}

func (c GenericContainer) CreateUser(user *system.OSUser) error {
  return user.CreateUser(path.Join(c.Location(), c.Name()))
}
