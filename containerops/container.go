package containerops

import(
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
