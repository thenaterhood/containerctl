package containerops

import(
  "fmt"
  "os"
  "path"
  "os/exec"
)

type ArchContainer struct {

  GenericContainer

}

func (c ArchContainer) Create() {

    c.GenericContainer.Create()

    if c.Installed() {
        fmt.Println(c.Name() + " is already installed")
        os.Exit(1)
    }
    dir := path.Join(c.Location(), c.Name())
    fmt.Println("Gonna pacstrap..." + dir)
    _, err := exec.Command("pacstrap", "-c", "-d", dir, "base", "--ignore", "linux").Output()

    if err != nil {
        fmt.Println(err)
    }
}
