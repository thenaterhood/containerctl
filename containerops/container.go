package containerops

import(
    "fmt"
    "os/exec"
    "os"
    "path"
    "path/filepath"
    "strings"
    "github.com/thenaterhood/containerctl/systemd"
)

type Container interface {
  New(name, location, uuid string, installed bool)
  Destroy()
  Create()
  Start()
  Stop()
  Exec(args ...string)

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

func (c GenericContainer) Create() {
    os.Mkdir(path.Join(c.Location(), c.Name()), 0700)
}

func (c GenericContainer) Destroy() {
    dir := path.Join(c.Location(), c.Name())
    os.RemoveAll(dir)
}

func (c GenericContainer) createMachineId() {
    dir := path.Join(c.Location(), c.Name())

    uuidbytes, _ := exec.Command("uuidgen").Output()
    uuid := string(uuidbytes[:37])
    machineid, _ := os.OpenFile(path.Join(dir, "etc", "machine-id"), os.O_WRONLY, 0600)
    machineid.WriteString(strings.Replace(uuid, "-", "", -1))
}

func (c GenericContainer) Exec(args ...string) {

    dir := path.Join(c.Location(), c.Name())
    command := append([]string{"-D", dir}, args...)

    out, err := exec.Command("systemd-nspawn", command...).Output()
    if err != nil {
      fmt.Printf("%s %s %s %s failed: %s\n, %s", "systemd-nspawn", "-D", dir, strings.Join(args, " "), err.Error(), out)
    }
}

func (c GenericContainer) Start() {
  if c.Installed() {
    systemd.StartService(c.Name() + "-container")
  } else {
    fmt.Println(c.Name() + " does not have a system installed.")
  }
}

func (c GenericContainer) Stop() {
  if c.Installed() {
    systemd.RunMachinectlCmd("poweroff", c.Name())
  } else {
    fmt.Println(c.Name() + " does not have a system installed.")
  }
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
