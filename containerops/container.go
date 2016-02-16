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

type Container struct {
    Name string
    Location string
    Installed bool
    Arch string
    Distro string
    Version string
    Uuid string
}

func (c Container) createMachineId() {
    dir := path.Join(c.Location, c.Name)

    uuidbytes, _ := exec.Command("uuidgen").Output()
    uuid := string(uuidbytes[:37])
    machineid, _ := os.OpenFile(path.Join(dir, "etc", "machine-id"), os.O_WRONLY, 0600)
    machineid.WriteString(strings.Replace(uuid, "-", "", -1))
}

func (c Container) Create() {
    os.Mkdir(path.Join(c.Location, c.Name), 0700)
}

func (c Container) Destroy() {
    dir := path.Join(c.Location, c.Name)
    os.RemoveAll(dir)
}

func (c Container) InstallArch() {
    if c.Installed {
        fmt.Println(c.Name + " is already installed with " + c.Distro)
        os.Exit(1)
    }

    dir := path.Join(c.Location, c.Name)
    _, err := exec.Command("pacstrap", "-c", "-d", dir, "base", "--ignore", "linux").Output()

    if err != nil {
        fmt.Println(err)
    }
}

func (c Container) InstallDebian() {
    if c.Installed {
        fmt.Println(c.Name + " is already installed with " + c.Distro)
        os.Exit(1)
    }

    dir := path.Join(c.Location, c.Name)
    _, err := exec.Command("debootstrap", "--arch="+c.Arch, c.Version, dir).Output()

    c.createMachineId()
    c.aptInstall("dbus")
    c.aptInstall("coreutils")
    c.runCommand("/lib/systemd/systemd-sysv-install", "enable", "dbus")

    if err != nil {
        fmt.Println(err)
    }
}

func (c Container) aptInstall(pkg string) {
    c.runCommand("apt-get", "install", "-y", "--allow-unauthenticated", pkg)
}

func (c Container) runCommand(args ...string) {
    dir := path.Join(c.Location, c.Name)
    command := append([]string{"-D", dir}, args...)

    out, err := exec.Command("systemd-nspawn", command...).Output()
    if err != nil {
      fmt.Printf("%s %s %s %s failed: %s\n, %s", "systemd-nspawn", "-D", dir, strings.Join(args, " "), err.Error(), out)
    }
}

func (c Container) Start() {
  if c.Installed {
    systemd.StartService(c.Name + "-container")
  } else {
    fmt.Println(c.Name + " does not have a system installed.")
  }
}

func (c Container) Stop() {
  if c.Installed {
    systemd.RunMachinectlCmd("poweroff", c.Name)
  } else {
    fmt.Println(c.Name + " does not have a system installed.")
  }
}

func Find(dir string) []*Container {

    var names []string
    var containers []*Container

    names, _ = filepath.Glob(dir + "/*")
    for _, name := range names {
        fi, _ := os.Stat(name)

        if fi.IsDir() {
            containers = append(containers, Load(name))
        }
    }
    return containers
}

func Load(dir string) *Container {

    c := new(Container)
    c.Location, c.Name = path.Split(dir)

    contents, _ := filepath.Glob(dir + "/*")

    if len(contents) == 0 {
        c.Installed = false

    } else {
        c.Installed = true

    }
    return c
}

func LoadMultiple(dir string, names []string) []*Container {

  var loaded []*Container

  for _, name := range names {
    loaded = append(loaded, Load(path.Join(dir, name)))
  }

  return loaded
}
