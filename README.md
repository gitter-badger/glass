# glassbox [![Build Status](https://travis-ci.org/acondolu/glassbox.svg?branch=master)](https://travis-ci.org/acondolu/glassbox)
There's no place like a safe cloud

## Description
Glass is a protocol for routing packets based on public-key cryptography.
It is also a framework for end-to-end-encrypted multi-application communication.

Currently being implemented in [Go](https://golang.org).

## Terminology
* **router**: the service provider, a public server routing packets to users and apps;
* **app instance**: the combination of user + app + device;
* **main app instance**: the instance of an application which acts as a server;
* **admin app**: the user application handling the configuration and the status of the service;

## App interface
Applications can use the App API to easily interface with the Glass protocol. New apps are initialized by authorization tokens, generated by the Admin App.

Take a look at [Go App API](https://github.com/acondolu/glassbox/wiki/Go-API) for more information and examples.

## Admin app
The main instance of the Admin app is the only application that can access and handle the user's private key; it is also the only mandatory component of the framework.

The Admin app handles the packets encrypted with user's public key, verifies them and redirects them to the recipient application and device.

Start up configuration:
* router's hostname
* router's public key
* user's private key


## Inspiring existing protocols
* XMPP
* SOCKS5
* Mobile IPv6
* HTTP Cloud API
