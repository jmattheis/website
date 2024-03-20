# jmattheis/website

This may seems like just a normal website, but it is much more. 
This server abuses various protocols that they will transmit data of my website.

Try one of the following commands:

```
curl     dict://jmattheis.de/show:server
curl      ftp://jmattheis.de
curl   gopher://jmattheis.de
curl     http://jmattheis.de
curl    https://jmattheis.de
curl    'imap://jmattheis.de/INBOX;UID=1' -u ':'
curl     pop3://jmattheis.de/1
dict -h         jmattheis.de -I
dig            @jmattheis.de +tcp +short
docker -H       jmattheis.de inspect -f '{{.Value}}' start
finger          jmattheis.de
gemget gemini://jmattheis.de/ -o-
netcat          jmattheis.de 23
redis-cli -h    jmattheis.de lrange start 0 0
ssh             jmattheis.de
telnet          jmattheis.de 23
whois -h        jmattheis.de .
wscat -c        jmattheis.de
```
