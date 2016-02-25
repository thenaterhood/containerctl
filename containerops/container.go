package containerops

import(
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "github.com/thenaterhood/containerctl/system"
)

type Container interface {
  Destroy() error
  Create() error
  Start() error
  Stop() error
  Exec(args ...string) error
  UpdateUser(*system.OSUser) error
  CreateUser(*system.OSUser) error

  Name() string
  Location() string
  Installed() bool
  Uuid() string
}

// Finds directory-based containers stored in a particular directory
// It accepts a string and returns []Container
// This function will only check that what it found is a directory, not
// that it is necessarily a valid container; it's assumed that it is.
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

// Loads a single container, given its full path.
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

// Loads multiple containers, given a path and a slice of container
// names. This function will load and return ONLY the containers requested.
func LoadMultiple(dir string, names []string) []Container {

  var loaded []Container

  for _, name := range names {
    loaded = append(loaded, Load(path.Join(dir, name)))
  }

  return loaded
}

// Converts any Container into a GenericContainer, which can be used
// in converting a container to another type
func ToGenericContainer(c Container) GenericContainer {
  var ctr GenericContainer
  ctr.location = c.Location()
  ctr.name = c.Name()
  ctr.uuid = c.Uuid()
  ctr.installed = c.Installed()

  return ctr
}

// Loads the machine id from the container (<container>/etc/machine-id)
// and returns it as a string.
func getContainerUuid(c Container) string {
  location := path.Join(c.Location(), c.Name(), "etc", "machine-id")
  contents, err := ioutil.ReadFile(location)
  var uuidstr string

  if err == nil {
    uuidstr = string(contents)
  }

  return uuidstr
}
