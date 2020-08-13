# Snorlax REST Client 

Snorlax is a REST client written in Go.

![Snorlax](/images/snorlax.jpg)

## Usage

### Creating a simple client.
```golang
client := snorlax.NewClient()
```

### Configuring the client using `ClientOptions`.
```golang
client := snorlax.NewClient(
	snorlax.WithBaseURL("https://www.example.com"),
	snorlax.WithRequestOptions(
		snorlax.WithHeader("Accept", "application/json"),
		snorlax.WithHeader("Content-Type", "application/json"),
	),
)
```

### Performing a simple request.
```golang
res, err := client.Get(context.Background(), "/example", nil)
if err != nil {
	log.Fatal(err)
}
```

### Performing a request with query parameters.
```golang
	params := make(url.Values)
	params.Set("name", "Snorlax")
	params.Set("number", 143")

	res, err := client.Get(context.Background(), "/example", params)
	if err != nil {
		log.Fatal(err)
	}
```

### Performing a request with a body.
```golang
payload := []byte("{\"name\": \"Snorlax\", \"number\": 143}")

res, err := client.Post(context.Background(), "/example", nil, bytes.NewBuffer(payload))
if err != nil {
	log.Fatal(err)
}
```

### Extracting JSON out of a response.
```golang
type Pokemon struct {
	Name 	string `json:"name"`
	Number 	string `json:"number"`
}

res, err := client.Get(context.Background(), "/example", nil)
if err != nil {
	log.Fatal(err)
}

var pokemon Pokemon
if err = res.JSON(pokemon); err != nil {
	log.Fatal(err)
}
```

## Contributing
Contributions are welcome! Feel free to submit pull requests. :) 
