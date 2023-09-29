# PAWS - Port Availability Waiting System

![paws-logo.jpg](assets/paws-logo.jpg)

Before you leap, let PAWS take a peep!

## What is PAWS?

PAWS can wait until multiple ports are read. It is a bit special, because it does not only can check if a port is available, but also if the expected service is running. It was created to aid in our integration and local testing.

## Beware: Work In Progress

If you stumbled on this page: We are actively working on this! Stuff may change and we can't promise that it will work for you or even work at all!

## Installation

```
go install github.com/metatexx/paws@latest 
```

## Usage examples

Checking if your localhost runs sshd

```
paws 22
```

* [Examples](examples.md)
* [Docker script example](examples/docker-mysql.sh)

You can also use the "cheat" command:

```
paws cheat --list
paws cheat examples
paws cheat docker
```

And even run the docker script from the cheats. You need docker on your system for that to work.

```
paws cheat docker 2>&1 | sh
```

## Usage

```
paws --help-long
```

```
usage: paws [<flags>] <command> [<args> ...]

Port Availability Waiting System: Before you leap, let PAWS take a peep.

# Use --help-long for more info and --cheats for examples

> Anonymous kitten: 'Curiosity checked the port!'

Flags:
      --help              Show context-sensitive help
      --version           Show application version.
  -h, --host="localhost"  Default host if none is given in the ports
  -q, --quiet             Don't output anything (only return code)
  -p, --progress          Output dots while scanning (one dot for each port check)
  -w, --wait=WAIT         Waiting time to wait for the ports to appear
  -t, --timeout=250ms     Timeout for all connect and read timeouts
  -d, --delay=250ms       Minimal time between checks when waiting

Args:
  (service:)(host:)port(-udp|-tcp)  Ports to scan (if you need host but without service name you can use use ':host:port'

Commands:
help [<command>...]
    Show help.


check [<flags>] (host:)port(-udp|-tcp)...
    checks the given port list

    -h, --host="localhost"  Default host if none is given in the ports
    -q, --quiet             Don't output anything (only return code)
    -p, --progress          Output dots while scanning (one dot for each port check)
    -w, --wait=WAIT         Waiting time to wait for the ports to appear
    -t, --timeout=250ms     Timeout for all connect and read timeouts
    -d, --delay=250ms       Minimal time between checks when waiting

cheat [<flags>] [<label>]
    Shows cheats for paws

    --list            List available cheats
    --save=DIRECTORY  Saves the cheats to the given directory
```

> *Anonymous kitten: Curiosity checked the port!*

---
*Copyright METATEXX GmbH 2023 / MIT-License*
