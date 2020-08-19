# Snorlax

[![Build](https://github.com/nickcorin/snorlax/workflows/Go/badge.svg?branch=master)](https://github.com/nickcorin/snorlax/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/nickcorin/snorlax?style=flat-square)](https://goreportcard.com/report/github.com/nickcorin/snorlax)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/nickcorin/snorlax)

Snorlax is a simple REST client written in Go.

![Snorlax](/images/snorlax.jpg)

## Installation

To install `snorlax`, use `go get`:
```
go get github.com/nickcorin/snorlax
```

Import the `snorlax` package into your code:
```golang
package main

import "github.com/nickcorin/snorlax"

func main() {
	client := snorlax.NewClient(nil)
}
```

## Usage

#### Creating a simple client.
```golang
client := snorlax.NewClient(nil)
```

#### Configuring the client using `ClientOptions` and `CallOption`s.
```golang
client := snorlax.NewClient(&snorlax.ClientOptions{
		BaseURL: 		"https://www.example.com",
		CallOptions: 	[]snorlax.CallOption{
			snorlax.WithHeader("Content-Type", "application/json"),
		},
	}
)
```

#### Performing a simple request.
```golang
res, err := client.Get(context.Background(), "/example", nil)
if err != nil {
	log.Fatal(err)
}
```

#### Performing a request with query parameters.
```golang
params := make(url.Values)
params.Set("name", "Snorlax")
params.Set("number", 143")

res, err := client.Get(context.Background(), "/example", params)
if err != nil {
	log.Fatal(err)
}
```

#### Performing a request with a body.
```golang
payload := []byte("{\"name\": \"Snorlax\", \"number\": 143}")

res, err := client.Post(context.Background(), "/example", nil, bytes.NewBuffer(payload))
if err != nil {
	log.Fatal(err)
}
```

#### Performing a request with `CallOption`s.
```golang
username, password := "testuser", "testpassword"

res, err := client.Get(context.Background(), "/example", nil, WithBasicAuth(username, password))
if err != nil {
	log.Fatal(err)
}
```

#### Extracting JSON out of a response.
```golang
type Pokemon struct {
	Name 	string `json:"name"`
	Number 	int    `json:"number"`
}

res, err := client.Get(context.Background(), "/example", nil)
if err != nil {
	log.Fatal(err)
}

var pokemon Pokemon
if err = res.JSON(&pokemon); err != nil {
	log.Fatal(err)
}
```

## Contributing
Please feel free to submit issues, fork the repositoy and send pull requests!

## License
This project is licensed under the terms of the MIT license.
