package containerops

import(
    "fmt"
    "os/exec"
    "os"
    "path"
    "path/filepath"
    "strings"
    "time"
    "github.com/thenaterhood/containerctl/systemd"
)

type Container interface {
  Destroy() error
  Create() error
  Start() error
  Stop() error
  Exec(args ...string) error

  Name() string
  Location() string
  Installed() bool
  Uuid() string
}

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
    err := os.Mkdir(path.Join(c.Location(), c.Name()), 0700)
    return err
}

func (c GenericContainer) Destroy() error {
    err := c.Stop()
    time.Sleep(100 * time.Millisecond)

    if err == nil {
      dir := path.Join(c.Location(), c.Name())
      err = os.RemoveAll(dir)
    }

    return err
}

func (c GenericContainer) createMachineId() {
    dir := path.Join(c.Location(), c.Name())

    uuidbytes, _ := exec.Command("uuidgen").Output()
    uuid := string(uuidbytes[:37])
    machineid, _ := os.OpenFile(path.Join(dir, "etc", "machine-id"), os.O_WRONLY, 0600)
    machineid.WriteString(strings.Replace(uuid, "-", "", -1))
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

func Find(dir string) []Container {

    var names []string
    var containers []Container

    names, _ = filepath.Glob(dir + "/*")
    for _, name := range names {
        fi, _ := os.Stat(name)

        if fi.IsDir() {
            containers = append(containers, Load(name))
        }
    }
    return containers
}

func Load(dir string) Container {

    location, name := path.Split(dir)

    c := new(GenericContainer)
    c.location = location
    c.name = name
    c.installed = false

    release_files, _ := filepath.Glob(dir+"/etc/*-release")

    if len(release_files) > 0 {
      c.installed = true
      release := release_files[0]

      _, err := os.Stat(dir+"/etc/machine-id")
      if err == nil {

      }

      switch release {
        case
        "os-release":

        deb := new(DebianContainer)

        deb.location = location
        deb.installed = true
        deb.name = name

        return deb
        break

        case
        "arch-release":

        arch := new(ArchContainer)

        arch.location = location
        arch.installed = true
        arch.name = name

        return arch

        break

      }
    }

    return c
}

func LoadMultiple(dir string, names []string) []Container {

  var loaded []Container

  for _, name := range names {
    loaded = append(loaded, Load(path.Join(dir, name)))
  }

  return loaded
}

func ToGenericContainer(c Container) GenericContainer {
  var ctr GenericContainer
  ctr.location = c.Location()
  ctr.name = c.Name()
  ctr.uuid = c.Uuid()
  ctr.installed = c.Installed()

  return ctr
}
