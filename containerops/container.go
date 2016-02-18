package containerops

import(
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
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
    installed := false
    release_files, _ := filepath.Glob(dir+"/etc/*-release")
    var release string
    var ctr Container

    if len(release_files) > 0 {
      installed = true
      release = release_files[0]
    }

    switch release {
      case
      "os-release":

      deb := new(DebianContainer)

      deb.location = location
      deb.installed = installed
      deb.name = name
      deb.uuid = getContainerUuid(deb)

      ctr = deb

      case
      "arch-release":

      arch := new(ArchContainer)

      arch.location = location
      arch.installed = installed
      arch.name = name
      arch.uuid = getContainerUuid(arch)

      ctr = arch

      default:
        gctr := new(GenericContainer)

        gctr.location = location
        gctr.installed = installed
        gctr.name = name
        gctr.uuid = getContainerUuid(gctr)

        ctr = gctr
    }

    return ctr
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

func getContainerUuid(c Container) string {
  location := path.Join(c.Location(), c.Name(), "etc", "machine-id")
  contents, err := ioutil.ReadFile(location)
  var uuidstr string

  if err == nil {
    uuidstr = string(contents)
  }

  return uuidstr
}
