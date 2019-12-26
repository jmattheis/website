+++
title = "Setup a forwarding DNS Sinkhole with DNS over TLS&HTTPS"
description = "Setup a forwarding DNS Sinkhole with DNS over TLS&HTTPS"
toc = true
type = "blog"
tags = [
    "tutorial",
    "dns"
]
date = "2019-12-26"
+++

## Install CoreDNS

CoreDNS is a DNS server written in Go.
It features an extensive plugin system for configuring it to your needs.
In my testing, CoreDNS just worked, so I didn’t try any other DNS server.

CoreDNS provides [pre-compiled binaries](https://github.com/coredns/coredns/releases/latest)
and [docker images](https://hub.docker.com/r/coredns/coredns/). 
You can also [build it from source](https://github.com/coredns/coredns#compilation-from-source).

For the simplicity of this tutorial, I’ll use the pre-compiled binary. As of the time of writing, the latest version is `1.6.6`.

Download the archive:
```bash
$ wget https://github.com/coredns/coredns/releases/download/v1.6.6/coredns_1.6.6_linux_amd64.tgz
```

Extract the archive:
```bash
$ tar -xvf coredns_1.6.6_linux_amd64.tgz
```

The archive contains an executable named `coredns`, which we will later need to start the server.

## Configure CoreDNS

CoreDNS will be configured via a file that can be defined via `-conf` (default: ./Corefile). We’ll go with the default file.

First, we configure a server block that matches anything,
because we want that server to handle all queries (whether to forward or block the domain).

Each server block starts with a zone and is followed by braces `{ .. }`. 
This is mostly irrelevant for us, as we only want to forward queries 
and not host a [Authoritative DNS server](https://en.wikipedia.org/wiki/Name_server#Authoritative_name_server) (A server which has the original zone records for a domain).

(Comments start with `#`)

-> File: `./Corefile`
```bash
# . matches everything
# :53 listen on port 53 (default DNS port)
.:53 { }
```

Now, we configure our first plugins:

-> File: `./Corefile`

```bash
.:53 {
  # the any plugin blocks any queries by responding with a short reply
  # See https://tools.ietf.org/html/rfc8482 for more information.
  any

  # log errors to standard out
  errors

  # for better verification, we add logging of all requests
  log
}
```

### Forward DNS

Next up, we configure CoreDNS to forward our queries to an existing DNS server.
I’ll use [cloudflare](https://developers.cloudflare.com/1.1.1.1/setting-up-1.1.1.1/) but you can choose the one you trust 
(see f.ex: [privacytools.io/providers/dns/](https://www.privacytools.io/providers/dns/)).

-> File: `./Corefile`

```bash
.:53 {
  any
  errors
  log

  # forward is the plugin name
  # the second parameter (.) is the base domain to match . = anything
  # the other parameter (before the {) are the endpoints to forward to.
  # tls:// means that DNS over TLS should be used for the communication
  # Also supported are:
  # dns:// -> normal unencrypted DNS
  # https:// -> DNS over HTTPS
  # grpc:// -> DNS over gRPC
  forward . tls://1.1.1.1 tls://1.0.0.1 {

    # the server name will be used in the TLS negotiation.
    tls_servername cloudflare-dns.com

    # the duration for checking the health of the upstream DNS server
    health_check 60s

  }
}
```

Let’s check our configuration. Start the core DNS server with:
```bash
$ sudo ./coredns
```
CoreDNS requires sudo because we use port 53, you can work around this by changing this to an unused port which is over 1000.

With a started server, we can make a test request with dig:

```dns
$ dig @localhost -p 53 google.com
; <<>> DiG 9.14.8 <<>> @localhost -p 53 google.com
; (2 servers found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 33196
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;google.com. IN A

;; ANSWER SECTION:
google.com. 181 IN A 172.217.17.78

;; Query time: 197 msec
;; SERVER: ::1#53(::1)
;; WHEN: Wed Dec 25 12:06:34 CET 2019
;; MSG SIZE  rcvd: 65
```

The server works, because it returned us the IP for google.com `172.217.17.78` (it may be different for you).


### Block Domains

Blocking domains can be done in different ways. IMO, the easiest way is by using host files.
There are many online resources, for host files with hosts that serve malware or ads.
In this tutorial, we use the basic version of [github.com/StevenBlack/hosts](https://github.com/StevenBlack/hosts).
But first, we try out how it works.

(Comments start with `#`)


-> File: `./hosts`
```hosts
# the host file has a really simple format.
# it starts with an ip address followed by one or more host names.
# In this example we resolve google.com to 0.0.0.0
0.0.0.0 google.com
```

0.0.0.0 is mostly used for blocking the domain.
See [github.com/StevenBlack/hosts](https://github.com/StevenBlack/hosts#we-recommend-using-0000-instead-of-127001).

> We prefer to use 0.0.0.0, which is defined as a non-routable meta-address used to designate an invalid, unknown, or non-applicable target.
>
> Using 0.0.0.0 is empirically faster, possibly because there’s no wait for a timeout resolution. It also does not interfere with a web server that may be running on the local PC.


-> File: `./Corefile`

```bash
.:53 {
  any
  errors
  log

  # hosts serve zone data from a hosts file
  hosts ./hosts {
    # fallthrough passes the request to the next plugin if it couldn’t
    # be found inside the hosts file. Without this, no domain could be 
    # resolved because the forward plugin will never be executed.
    fallthrough
  }

  forward . tls://1.1.1.1 tls://1.0.0.1 {
    tls_servername cloudflare-dns.com
    health_check 60s
  }
}
```

> The order of the plugins inside the Corefile doesn’t matter, 
> you could add the hosts plugin after the forward plugin and it would still have the same behavior.
> The ordering of the plugins is defined in https://github.com/coredns/coredns/blob/master/plugin.cfg


Start the core DNS server:
```bash
$ sudo ./coredns
```

Check if google.com can be resolved.
```dns
$ dig @localhost -p 53 google.com
[removed bloat]

;; ANSWER SECTION:
google.com. 3600 IN A 0.0.0.0

[removed bloat]
```

google.com resolves to `0.0.0.0` which blocks the domain.
As we do not want to block google.com but use the host file from StevenBlack/hosts,
we remove the old hosts file
```bash
$ rm hosts
```
and download the hosts file from the GitHub repository:
```bash
$ wget https://github.com/StevenBlack/hosts/raw/master/hosts
```

After a restart of the CoreDNS server, some malicious domains will be blocked and google.com is available again.

### DNS over TLS (DoT) & DNS over HTTPS (DoH)


For DoT/DoH to work correctly you need a domain with a valid TLS certificate,
you can get one via certibot https://certbot.eff.org/ or purchase one.
You also have to create an entry in your domain settings to point your domain (or a subdomain)
to the server where CoreDNS is hosted on.

In this tutorial we have our certificats at `/var/certs/full.pem` and `/var/certs/key.pem`.

```bash
# add tls://.:953 to listen for DoT connections
# add https://.:443 to listen for DoH connections
.:53 tls://.:953 https://.:443 {

  # add the TLS plugin with the certs
  tls /var/certs/full.pem /var/certs/key.pem

  any
  errors
  log
  hosts ./hosts {
    fallthrough
  }
  forward . tls://1.1.1.1 tls://1.0.0.1 {
    tls_servername cloudflare-dns.com
    health_check 60s
  }
}
```

Restart the CoreDNS server and now it serves DoT/DoH :D.


#### Use DoT in Android

* Open Settings
* Click on `Network & Internet`
* Click on `Advanced`
* Click on `Private DNS`
* Enter your domain or subdomain inside `Private DNS provider hostname`.

A log entry should appear if you open a website on your phone. I visited `jmattheis.de` and
it created this log entry:

```
[INFO] [redacted-ip]:41938 - 0 “A IN jmattheis.de. tcp 128 true 65535” NOERROR qr,rd,ra 153 0.006815507s
```

### Cache results

Caching results will reduce the traffic on the upstream DNS server.

```bash
.:53 tls://.:953 https://.:443 {
  tls /var/certs/full.pem /var/certs/key.pem
  any
  errors
  log
  hosts ./hosts {
    fallthrough
  }
  forward . tls://1.1.1.1 tls://1.0.0.1 {
    tls_servername cloudflare-dns.com
    health_check 60s
  }

  # Cache for 60 seconds
  cache 60
}
```


