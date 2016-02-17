containerctl
============

Containerctl is a utility written in Go for working with systemd-nspawn containers. It is not feature complete (you will still need to use machinectl and other manual commands) but is capable of creating, installing, and destroying Debian and ArchLinux containers.

**This software is not production ready and its API is subject to change until
the release of v1.0.0+**

Requirements
------------
* systemd
* pacstrap (for creating Arch containers)
* debootstrap (for creating Debian containers)

Usage
------------

The CLI is NOT finalized, but currently looks like:

```
# Create and install an ArchLinux container named YourContainer
$ containerctl create-arch YourContainer

# Create and install a Debian container named YourContainer
$ containerctl create-debian YourContainer

# Destroy and delete a container (does NOT power it down first)
$ containerctl destroy YourContainer
```

License
-----------
Containerctl is licensed under the MIT license. The full license text can be found in the LICENSE file.

If you find containerctl useful, use it regularly, or build something cool around it, please consider contributing, providing feedback or simply dropping a line to say that containerctl is useful to you. Feedback from users is what keeps open source projects strong.
