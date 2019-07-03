# go-deadmanssnitch [![GoDoc](https://godoc.org/github.com/PremiereGlobal/go-deadmanssnitch?status.png)](http://godoc.org/github.com/PremiereGlobal/go-deadmanssnitch) [![Go Report Card](https://goreportcard.com/badge/github.com/PremiereGlobal/go-deadmanssnitch)](https://goreportcard.com/report/github.com/PremiereGlobal/go-deadmanssnitch) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/PremiereGlobal/go-deadmanssnitch/blob/master/LICENSE)
go-deadmanssnitch is a [Go](http://golang.org/) client library for the [Dead Man's Snitch](https://deadmanssnitch.com/docs/api/v1) API.

## Installation
Make sure you have a working Go environment. To install, simply run:

```
go get github.com/PremiereGlobal/go-deadmanssnitch
```

## Usage

```
package main

import (
  "github.com/PremiereGlobal/go-deadmanssnitch"
)

var apiKey = "" // Set your api key here

func main() {

  client := deadmanssnitch.NewClient(apiKey)
  ...
}
```

For more information, read the [godoc package documentation](http://godoc.org/github.com/PremiereGlobal/go-deadmanssnitch).

## Examples

### Check-In

```
  var	snitchToken = "" // Set your snitch token here
  err := client.CheckIn(snitchToken)
  if err != nil {
    panic(err)
  }
```

### List All Snitches

```
  snitches, err := client.ListSnitches([]string{})
  if err != nil {
    panic(err)
  }
```

### Create Snitch

```
  snitch := deadmanssnitch.Snitch {
    Name:      "testSnitch",
    Interval:  "hourly",
    AlertType: "basic",
    Tags:      []string{"test"},
    Notes:     "This is an example snitch",
  }

  createdSnitch, err := client.CreateSnitch(&snitch)
  if err != nil {
    panic(createdSnitch)
  }
```

## Testing the Client Library
Tests will validate API calls by creating a test snitch, checking in and updating the snitch using all of the methods.  It will then delete the snitch after waiting `wait` seconds (to allow for manual verification).

```
go test -v --args -apikey=<apiKey> -wait 30
```
