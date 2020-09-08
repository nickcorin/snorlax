<p align="center">
<h1 align="center">Snorlax</h1>
<p align="center">Simple HTTP and REST client library written in Go</p>
</p>
<p align="center">
<p align="center"><a href="https://github.com/nickcorin/snorlax/actions?query=workflow%3AGo"><img src="https://github.com/nickcorin/snorlax/workflows/Go/badge.svg?branch=master" alt="Build Status"></a> <a href="https://goreportcard.com/report/github.com/nickcorin/snorlax"><img src="https://goreportcard.com/badge/github.com/nickcorin/snorlax?style=flat-square" alt="Go Report Card"></a> <a href="http://godoc.org/github.com/nickcorin/snorlax"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="GoDoc"></a> <a href="LICENSE"><img src="https://img.shields.io/github/license/nickcorin/snorlax" alt="License"></a></p>
</p>
<p align="center">
<img src="/images/snorlax.jpg" />
</p>

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
	client := snorlax.DefaultClient
}
```

## Usage

#### Using the DefaultClient.
```golang
client := snorlax.DefaultClient
```

#### Configuring the client using `ClientOptions`.
```golang
client := snorlax.NewClient(snorlax.ClientOptions{
		BaseURL: "https://www.example.com",
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

#### Performing a request with `PreRequestHook`s.
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
