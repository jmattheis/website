+++
title = "Website from Jannis Mattheis."
description = "Website from Jannis Mattheis."
+++

Hey there! I'm Jannis Mattheis, a developer from Germany.

This may seem like a normal website, but it is much more.
This server abuses various protocols to transfer content of my website.

Currently supported are: 
  dict, dns(tcp), ftp, gopher, http/https, imap, pop3,  ssh, telnet/tcp, websocket and whois

Try one of the following commands:

```
curl   dict://jmattheis.de/show:server
curl    ftp://jmattheis.de
curl gopher://jmattheis.de
curl   http://jmattheis.de
curl  https://jmattheis.de
curl  "imap://jmattheis.de/INBOX;UID=1" -u ":"
curl   pop3://jmattheis.de/1
dict -h       jmattheis.de -I
dig          @jmattheis.de +tcp +short
docker -H     jmattheis.de inspect start -f '{{.Value}}'
netcat        jmattheis.de 23
ssh           jmattheis.de
telnet        jmattheis.de 23
whois -h      jmattheis.de .
wscat -c      jmattheis.de
```