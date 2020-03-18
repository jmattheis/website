+++
title = "Website from Jannis Mattheis."
description = "Website from Jannis Mattheis."
+++

Hey there! I'm Jannis Mattheis, a developer from Germany.

This may seem like a normal website, but it is much more.
This server abuses various protocols to transfer contents from my website.

Currently supported are: 
  http/https, websocket, telnet/tcp, whois, dns(tcp), ftp and ssh

Try one of the following commands:

```
curl http://jmattheis.de
curl  ftp://jmattheis.de
dig        @jmattheis.de +tcp +short
netcat      jmattheis.de 23
ssh         jmattheis.de
telnet      jmattheis.de 23
whois -h    jmattheis.de .
wscat -c    jmattheis.de
```