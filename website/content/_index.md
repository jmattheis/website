+++
title = "Website from Jannis Mattheis."
description = "Website from Jannis Mattheis."
+++

Hey there! I'm Jannis Mattheis, a developer from Germany.

This may seem like a normal website, but it is much more.
This server abuses various protocols to transfer content of my website.

Currently supported are: 
  dns(tcp), ftp, http/https, imap, pop3,  ssh, telnet/tcp, websocket, whois

Try one of the following commands:

```
curl   ftp://jmattheis.de
curl  http://jmattheis.de
curl https://jmattheis.de
curl "imap://jmattheis.de/INBOX;UID=1" -u ":"
curl  pop3://jmattheis.de/1
dig         @jmattheis.de +tcp +short
netcat       jmattheis.de 23
ssh          jmattheis.de
telnet       jmattheis.de 23
whois -h     jmattheis.de .
wscat -c     jmattheis.de
```