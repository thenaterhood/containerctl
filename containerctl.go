package main

import(
    "fmt"
    "os"
    "github.com/thenaterhood/containerctl/containerops"
)

func main() {

    container_path := "/var/lib/container"

    if len(os.Args) < 2 {
        fmt.Println("Whoops! Too few arguments.")
        os.Exit(1)
    }

    action := os.Args[1]

    if len(os.Args) < 3 {
        switch action {
            case
            "list":
            for _, c := range containerops.Find(container_path) {
                fmt.Println(c.Name())
            }
        }
        os.Exit(0)
    }

    on_containers := containerops.LoadMultiple(container_path, os.Args[2:])

    switch action {
        case
        "create",
        "make":
        for _, c := range on_containers {
          c.Create()
        }
        break

        case
        "create-arch":
        for _, c := range on_containers {
          fmt.Println("Installing archlinux into " + c.Name() + "...")
          gc := containerops.ToGenericContainer(c)
          ctr := containerops.ArchContainer{gc}
          ctr.Create()
        }
        break

        case
        "create-debian":
        for _, c := range on_containers {
          fmt.Println("Installing debian sid into " + c.Name() + "...")
          gc := containerops.ToGenericContainer(c)
          ctr := containerops.DebianContainer{gc}
          ctr.Create()
        }
        break

        case
        "destroy",
        "remove":
        for _, c := range on_containers {
          fmt.Println("Destroying " + c.Name())
          c.Destroy()
        }
        break

        case
        "poweron",
        "start":
        for _, c := range on_containers {
          c.Start()
        }
        break

        case
        "stop",
        "poweroff":
        for _, c := range on_containers {
          c.Stop()
        }
        break
    }
}
