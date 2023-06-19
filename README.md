# EnvLoader
## Description
EnvLoader is a simple library that allows you to load environment variables into a model struct. And it supports nested structs too!

## Installation
```bash
go get github.com/metinorak/envloader
```

## Usage
```go
package main

import (
    "fmt"
    "github.com/metinorak/envloader"
)

type Config struct {
    Database struct {
        Host     string `env:"DB_HOST"`
        Port     int    `env:"DB_PORT"`
        Username string `env:"DB_USERNAME"`
        Password string `env:"DB_PASSWORD"`
        Name     string `env:"DB_NAME"`
    }
    Server struct {
        Host string `env:"SERVER_HOST"`
        Port int    `env:"SERVER_PORT"`
    }
}

func main() {
    // Example environment variables
    // DATABASE.HOST=localhost
    // DATABASE.PORT=3306
    // DATABASE.USERNAME=root
    // DATABASE.PASSWORD=secret
    // DATABASE.NAME=example
    // SERVER.HOST=localhost
    // SERVER.PORT=8080

    // Following lines will load environment variables into Config struct

    var config Config
    envLoader := envloader.New()

    err := envLoader.Load(&config)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", config)
}
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
```
```
