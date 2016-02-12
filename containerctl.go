package main

import(
        "fmt"
        "path/filepath"
        "os/exec"
        "os"
)

func find_containers(paths []string) []string {

        var containers, dirs []string

        for _, path := range paths {
                dirs, _ = filepath.Glob(path + "/*")
                containers = append(containers, dirs...)

        }
        return containers
}

func run_machinectl_cmd(cmd, container string) {
        _, err := exec.Command("machinectl " + cmd + " " + container).Output()

        if err != nil {
                fmt.Println(cmd + " failed on " + container)
                fmt.Println(err)
        }
}

func poweroff_containers(containers []string) {

        for _, container := range containers {
                run_machinectl_cmd("poweroff", container)
        }
}

func poweron_containers(containers []string) {
        for _, container := range containers {
                run_machinectl_cmd("start", container)
        }
}

func strInSlice(s string, slice []string) bool {
        for _, str := range slice {
                if s == str {
                        return true
                }
        }
        return false
}

func main() {
        paths := []string{"/var/lib/container"}

        if len(os.Args) < 3 {
                fmt.Println("Whoops! Too few arguments.")
                os.Exit(1)
        }

        action := os.Args[1]
        on_containers := os.Args[2:]

        containers := find_containers(paths)

        for _, rq_container := range on_containers {
                if ! strInSlice(rq_container, containers) {
                        fmt.Println(rq_container + " is not a known container.")
                        os.Exit(1)
                }
        }

        switch action {
                case
                "poweron",
                "start":
                poweron_containers(on_containers)
                break

                case
                "stop",
                "poweroff":
                poweroff_containers(on_containers)
                break
        }
}
