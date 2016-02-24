package main

import(
    "fmt"
    "os"
    "os/exec"
    "github.com/thenaterhood/containerctl/containerops"
    "github.com/thenaterhood/containerctl/system"
)

func main() {

    container_path := "/var/lib/container"

    if len(os.Args) < 2 {
        fmt.Println("Whoops! Too few arguments.")
        os.Exit(1)
    }

    action := os.Args[1]

    ctr_start_index := 2

    if action == "copy-user" {
      ctr_start_index += 1
    }

    on_containers := containerops.LoadMultiple(container_path, os.Args[ctr_start_index:])

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

        case
        "copy-user":
        known_users := system.LoadUsers("/")
        user := known_users.Find(os.Args[2])
        if user == nil {
          fmt.Println("User does not exist on host")
          break
        }

        for _, c := range on_containers {
          err := c.UpdateUser(user)
          if err != nil {
            fmt.Println(err)
          }
        }


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

        case
        "poweron",
        "start":
        for _, c := range on_containers {
          err := c.Start()
          if err != nil {
            fmt.Println(err)
          }
        }

        case
        "stop",
        "poweroff":
        for _, c := range on_containers {
          err := c.Stop()
          if err != nil {
            fmt.Println(err)
          }
        }

        default:
          out, _ := exec.Command("machinectl", os.Args[1:]...).Output()
          fmt.Println(string(out[:]))
    }
}
