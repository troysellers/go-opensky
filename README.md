# Simple Open Sky in Go

This is something I was playing with while starting to learn the Go programming language.
[OpenSky](https://opensky-network.org/) is an open source community based network that collects air traffic surveillance data. 

The program will run on the commandline and simply print the []state to stdout. 

```
go run opensky.go
```

if you have credentials for OpenSky and want to use them
```
go run opensky.go username password
```

As this only calls one method on the [API](https://opensky-network.org/apidoc/), the [/states/all](https://opensky-network.org/apidoc/rest.html#all-state-vectors) you probably won't see much benefit from authenticating.

 