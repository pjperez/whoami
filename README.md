# whoami
A simple DNS reflector with a twist!

## Introduction

Websites that use Global Load Balancing with multiple front-ends around the world need a way for you to reach them on your geographically closest endpoint. One way of doing that is resolving your query of www.example.com to the IP address of the endpoint in the closest datacenter. 

How do they know which one is the closest?

The idea is to geolocate the IP address of the DNS forwarder that sent them the request. In a common case, it will be the IP address of your ISP's DNS forwarder.

Now my question for you is: Are you sure your ISP's DNS forwarder is geolocated in the right region? Otherwise, you might be having a sub-par browsing experience. 

This also applies for the workplace, where multinational corporations might centralise their DNS services in one of their locations, so the rest of the branches might inadvertely be affected by this issue when consuming highly available web applications. Given that today's web applications are a big part of many people's workflow, why don't run a few checks to ensure everything is as it should be?

## How is a DNS reflector useful?

A DNS reflector is just a DNS server that returns the client's IP address when queried for a specific hostname, instead of a fixed IP address for that domain.

The client in that case would be your ISP's DNS forwarder, because when you try to resolve a hostname, this is roughly how it works:

Computer requests resolution for www.example.com to ISP's DNS Server -- ISP's DNS Server forwards (hence DNS forwarder) the request to the root servers, then based on the root servers reply the forwarder will send a new request (as a client!) to the authoritative name servers for that domain. These authoritative name servers are the ones responsible to return the IP address of the requested hostname (www.example.com) and at the same time will check on the client's IP address to decide on its geolocation and return the IP address closer to the client.

## What is this about?

This is the source code for the DNS reflector I'm running at whoami.fluffcomputing.com

Would you like to see what it does? Try resolving whoami.fluffcomputing.com

    pjperez@whoami:~$ dig +short whoami.fluffcomputing.com a
    81.139.58.34

That IP address is actually the IP address of your ISP DNS forwarder. In my case, my provider is BT (however they seem to outsource the service to one of those companies that intercepts your DNS requests to non-existing domains).

It gets better, why don't you try a txt record?

    pjperez@whoami:~$ dig +short whoami.fluffcomputing.com **txt**
    "London, United Kingdom"
  
What's that? That's the geolocation of your ISP DNS forwarder! Based on Maxmind's [GeoLite2](http://dev.maxmind.com/geoip/geoip2/geolite2/) DB.

### It also works nicely from Windows!

    PS D:\whoami> (Resolve-DnsName whoami.fluffcomputing.com -Type txt).strings
    London, United Kingdom
    PS D:\whoami> (Resolve-DnsName whoami.fluffcomputing.com -Type a).IPAddress
    81.139.58.35
  
## Credits

Awesome geoip library to consume MaxMind's DB: https://github.com/oschwald/geoip2-golang

Awesome DNS library: https://github.com/miekg/dns - https://miek.nl/2014/August/16/go-dns-package/

## Disclaimer

This is part of my own exercises to learn Golang, please take it as what it is.

The accuracy of this software and the service on whoami.fluffcomputing.com are based on [MaxMind's geolite2 DB](http://dev.maxmind.com/geoip/geoip2/geolite2/).

My running copy of the DB on whoami.fluffcomputing.com was updated on the 3rd of October 2016.
