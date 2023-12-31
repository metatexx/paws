# Examples for the usage of PAWS

## Port & Service

Checking if your localhost runs sshd

```
paws 22
```

Multiple ways to check for dns available on ip 192.168.1.1

```
# checks for tcp-port open only
paws :192.168.1.1:53
paws routerdns:192.168.1.1:53
paws routerdns://192.168.1.1:53

# checks for dns server protocol on tcp-port
paws dns:192.168.1.1:53
paws dns://192.168.1.1:53
paws dns-tcp:192.168.1.1:53
paws dns-tcp://192.168.1.1:53

# checks for dns on udp-port 
paws dns-udp:192.168.1.1:53
paws dns-udp://192.168.1.1:53
```

Checking if your router has DNS, localhost has ssh and the gmail smtp is reachable. Using the `--host` (`-h`) flag to specify another default host.

```
paws -h 192.168.1.1 dns:53 ssh:localhost:22 :smtp.gmail.com:465
```

Checking if port 4222 replies with a string starting with "INFO". A nats-server does that. The schema name is not important here. But it needs to be one that is not supported with other checks.

```
paws replies://:4222?INFO
```

Checking if a mssql server is up and ready for queries

Notice: "ping failed" usually means that the authorisation data did not match!

```
paws 'mssql://user:pass@localhost:1433?db=master&q=SELECT 1'
```

A docker example can be found with

```
paws cheat docker
```

To run the example you need docker installed and do

```
paws cheat docker 2>&1 | sh
```

## (Log-)Files & Stdout

Waiting for a docker logs to be idle for 3 seconds and fails after 30 seconds

```
docker logs -f <container> | paws log -t 30s -i 3s 
```

The same as above but showing the data while waiting

```
docker logs -f <container> | paws log -t 30s -i 3s -p
```

Waiting for up to five minutes till the words 'ok' and 'done' appear in the given file.

```
paws -f /tmp/file.log -t 5m -S 'ok' -S 'done'
```

Fail if in the next 5 seconds the word 'failure' appears in the given file and show the line that fails

```
paws -f /tmp/file.log -t 5m -F 'failure' -P
```
