---
title: "Send Notifications to Android via REST-API"
url: ["blog/send-notifications-to-android-via-rest-api", "blog/2019/08/send-notifications-to-android-via-rest-api" ]
description: "A tutorial on how to send push notifications to your phone with Gotify."
date: "2019-08"
---

A while ago, I wanted to get notifications when 
certain events occur on my servers. 
This could be a SSH-Login or finishing a backup. 

Back then, I started to self-host services to be less dependent 
on third-parties and to gain control over my data. 
After some research, I couldn't find any maintained open source project 
which had the functionality I wanted, 
so I started my own named [Gotify](https://gotify.net/).
Gotify is a pretty simple service written in Go that exposes a WebSocket endpoint 
which can be used to subscribe to newly posted messages. 
Gotify does not depend on **any** third-party service to function, 
thus does not use google services to deliver push notifications to your phone.

In this blog post, I'll show you how to set up Gotify and send some messages to it.

## Setting Up Gotify

Gotify can be started via binary or Docker. In this tutorial, 
I'll use the provided Docker images as they are pretty easy to set up.

[hub.docker.com/r/gotify/server](https://hub.docker.com/r/gotify/server)

```bash
$ docker run -p 8080:80 \
             -v /var/gotify/data:/app/data \
             -e GOTIFY_DEFAULTUSER_PASS=secret \
             gotify/server:2.0.6
```
or as via docker-compose:
```yml
version: '3'
 
services:
  gotify:
    image: gotify/server:2.0.6
    ports:
      - 8080:80
    volumes:
      - "/var/gotify/data:/app/data"
    environment:
      GOTIFY_DEFAULTUSER_PASS: "secret"
```
(start with `docker-compose up -d`)

By default, Gotify uses SQLite as database. 
Thus, further simplifying the setup because no 
separate database is needed. SQLite should work well with a small user base, 
however if you have many concurrent users a different 
database may improve performance. 
Besides SQLite Gotify supports PostgreSQL and MySQL/MariaDB.

`/app/data` contains the database file (if SQLite is used), 
images for applications and other stuff. 
In this example the directory is mounted to `/var/gotify/data` 
this directory should be included in a backup.

`-e GOTIFY_DEFAULTUSER_PASS=secret` changes the password 
of the default user which will be created at startup.

Have a look at [gotify.net/docs/config](https://gotify.net/docs/config) 
for all configuration options (like different database settings).

## First Login / Definitions

By default, the default username/password is `admin`, 
however in this tutorial we changed the password to `secret`. 
With these credentials it's now possible 
to login into the WebUI at http://localhost:8080/ 
(use the port you specified while starting the docker container).

In the UI you can configure different things.

**Clients**: A client is a device or application 
that can manage other clients, messages and applications. 
However, a client is not allowed to send messages.

In this case your browser would be a client.

**Applications**: An application is a device or 
application that only can send messages.

An application could be a raspberry pi 
which notifies when it reboots.

## Sending a message

You need an application to send messages to Gotify. 
Only the user who created the application 
is able to see its messages. An application can be added via:

* WebUI: click the `apps`-tab in the upper right corner when logged in and add an application
* REST-API: `curl -u admin:secret -X POST https://yourdomain.com/application -F "name=test" -F "description=tutorial"` 
  See [API-Docs](https://gotify.github.io/api-docs/)

To authenticate as an application, you need the application token. 
The token is returned in the REST request and is viewable in the WebUI.

After copying the token you can simply use curl, 
HTTPie or any other http-client to push messages.

```
$ curl -X POST "http://localhost/message?token=<apptoken>" -F "title=my title" -F "message=my message" -F "priority=5"
$ http -f POST "http://localhost:8080/message?token=<apptoken>" title="my title" message="my message" priority="5"
```
Replace `<apptoken>` with your application token, 
it should look like this: `AKTlZf.InA3uZHK`.

`priority` currently only has an effect in the android app, 
0 = not intrusive 10 = very intrusive.

The UI will render the message as plain text, 
it is possible to render it as markdown with 
[extras](https://gotify.net/docs/msgextras).

You can use [gotify/cli](https://github.com/gotify/cli) 
to push messages. The CLI stores url and token in a config file.

```bash
$ gotify push -t "my title" -p 10 "my message"
$ echo my message | gotify push
```
[Install gotify/cli](https://github.com/gotify/cli).

## Android App

While the WebUI already creates notifications on new messages, 
Gotify also has an android app named 
[gotify/android](https://github.com/gotify/android). 
It is available in the 
[Play Store](https://play.google.com/store/apps/details?id=com.github.gotify), 
on [F-Droid](https://f-droid.org/de/packages/com.github.gotify/) 
and you can download the apk 
[on the releases page](https://github.com/gotify/android/releases/latest). 
Setup is straight forward, enter the url to your Gotify instance and login. 

Be aware: By default Android kills long-running apps as they drain the battery. 
With enabled battery optimization, Gotify will be killed, 
and you won't receive any notifications. 

Here is one way to disable battery optimization for Gotify.

* Open "Settings"
* Search for "Battery Optimization"
* Find "Gotify" and disable battery optimization

---

... and that's it! Thanks for reading my first blog post (:.
