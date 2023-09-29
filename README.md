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

Checking if your router has DNS, localhost has ssh and the gmail smtp is reachable

```
paws -h 192.168.1.1 dns:53 ssh:localhost:22 :smtp.gmail.com:465
```

Checking for some ports till they are up after running docker containers

```
#!/bin/sh
docker kill mariadb33307 1>/dev/null 2>&1
docker kill mysql33306 1>/dev/null 2>&1
set -e
echo "running mariadb in docker"
docker run -d --name=mariadb33301 -p 33301:3306 --rm --env MARIADB_USER=user --env MARIADB_PASSWORD=user --env MARIADB_ROOT_PASSWORD=root mariadb:latest
echo "running mysql in docker"
docker run -d --name=mysql33302 -p 33302:3306 --rm mysql/mysql-server:latest
echo "waiting for all ports to be up (noisy version, use -p or -q for more silent operation)"

paws -w 10s mariadb:33301 mysql:33302
RETURN=$?
if [ $RETURN -eq 0 ];
then
  echo "all ports were found"
  exit 0
else
  echo "some or all ports were missing"
fi

echo "removing docker containers"
docker kill mariadb33301
docker kill mysql33302


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
