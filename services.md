# Services that are supported directly or indirectly 

Not only can PAWS check for open ports, it can also examine numerous services more thoroughly.
 
- nats (port 4222)
- mysql (port 3306)
- mariadb (port 3306)
- mssql (port 1433)
- dhs tcp & udp (port 53)
- ssh (port 22)
- smtp (port 25)
 
- planned: imap-tls (port 993)
- planned: smtp-tls (port 465)

Another way to test is by using "reply://" with an query parameter:

- paws reply://imap.ionos.de:143?"* OK ["
