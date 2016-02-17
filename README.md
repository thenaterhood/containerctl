containerctl
============

Containerctl is a utility written in Go for working with systemd-nspawn containers. It is not feature complete (you will still need to use machinectl and other manual commands) but is capable of creating, installing, and destroying Debian and ArchLinux containers.

Usage
------------
The CLI is NOT finalized, but currently looks something like:

```
# Create and install an ArchLinux container named YourContainer
$ containerctl create-arch YourContainer

# Create and install a Debian container named YourContainer
$ containerctl create-debian YourContainer

# Destroy and delete a container (does NOT power it down first)
$ containerctl destroy YourContainer
```
