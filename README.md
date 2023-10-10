# PAWS - Port Availability Waiting System

![paws-logo.jpg](assets/paws-logo.jpg)

Before you leap, let PAWS take a peep!

## What is PAWS?

PAWS is a universal tool to wait for conditions being meet. It was all out ports and services at first but now can
also wait on (log-files). It was created to aid in our integration and local testing.

It has some unique features:

- It can check for ports being available
  - It can check for multiple known services (nats/mysql/mssql/ssh and more).
  - It can also check for a service to return some given text 
- It can check a file or stdout (and tails it automatically)
  - For multiple patterns to appear and succeeds after all of them are found.
  - Fail if any of some given patterns appear.
  - Fail when the file is not fully read or a condition was found after some time
  - Succeed if a file does not get output append after some time (being idle). This can be used to wait for services to stabilise.
  - Found patterns can be highlighted
  - Output can be transparently passed through.
  - It can be made fully quiet and only returns the return code to the shell
- It has multiple ways to handle timeouts.

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
* [Services](sevices.md)
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
