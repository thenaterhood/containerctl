package main

import(
    "fmt"
    "os"
    "os/exec"
    "github.com/thenaterhood/containerctl/containerops"
)

func main() {

    container_path := "/var/lib/container"

    if len(os.Args) < 2 {
        fmt.Println("Whoops! Too few arguments.")
        os.Exit(1)
    }

    action := os.Args[1]

    on_containers := containerops.LoadMultiple(container_path, os.Args[2:])

    switch action {
        case
        "create",
        "make":
        for _, c := range on_containers {
          err := c.Create()
          if err != nil {
            fmt.Println(err)
          }
        }
        break

        case
        "create-arch":
        for _, c := range on_containers {
          fmt.Println("Installing archlinux into " + c.Name() + "...")
          gc := containerops.ToGenericContainer(c)
          ctr := containerops.ArchContainer{gc}
          err := ctr.Create()
          if err != nil {
            fmt.Println(err)
          }
        }
        break

        case
        "create-debian":
        for _, c := range on_containers {
          fmt.Println("Installing debian sid into " + c.Name() + "...")
          gc := containerops.ToGenericContainer(c)
          ctr := containerops.DebianContainer{gc}
          err := ctr.Create()
          if err != nil {
            fmt.Println(err)
          }
        }
        break

        case
        "destroy",
        "remove":
        for _, c := range on_containers {
          fmt.Println("Destroying " + c.Name())
          err := c.Destroy()
          if err != nil {
            fmt.Println(err)
          }
        }
        break

        case
        "poweron",
        "start":
        for _, c := range on_containers {
          err := c.Start()
          if err != nil {
            fmt.Println(err)
          }
        }
        break

        case
        "stop",
        "poweroff":
        for _, c := range on_containers {
          err := c.Stop()
          if err != nil {
            fmt.Println(err)
          }
        }
        break

        default:
          out, _ := exec.Command("machinectl", os.Args[1:]...).Output()
          fmt.Println(string(out[:]))
    }
}
