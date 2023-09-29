#!/bin/sh
# can be run with `paws cheat docker 2>&1 | bash`
docker kill mariadb33301 1>/dev/null 2>&1
docker kill mysql33302 1>/dev/null 2>&1
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
