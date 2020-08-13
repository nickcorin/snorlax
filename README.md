# Snorlax REST Client 

![Snorlax](/images/snorlax.png)

## Usage

```golang
package main

import "github.com/nickcorin/snorlax"

func main() {
	// Creating a client, with some configuration options.
	client := snorlax.NewClient(
		snorlax.WithBaseURL("https://www.example.com),
		snorlax.WithRequestOptions(
			snorlax.WithHeader("Accept", "application/json"),
			snorlax.WithHeader("Content-Type", "application/json"),
		),
	)

	// Perform a simple request.
	res, err := client.Get(context.Background(), "/example", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Perform a request with query parameters.
	params := make(url.Values)
	params.Set("name", "Snorlax")
	params.Set("number", 143")

	res, err := client.Get(context.Background(), "/example", params)
	if err != nil {
		log.Fatal(err)
	}

	// Perform a request with body.
	payload := []byte("{\"name\": \"Snorlax\", \"number\": 143}")

	res, err := client.Post(context.Background(), "/example", nil,
		bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
}
```
