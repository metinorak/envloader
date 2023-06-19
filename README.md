# EnvLoader
## Description
EnvLoader is a simple library that allows you to load environment variables into a model struct. And it supports nested structs too!

## Features
- Supports nested structs
- Field names are converted to upper snake case by default
- But you can specify custom field names with `env` tags

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

// An example nested struct
type Config struct {
    Database struct {
        Host     string
        Port     int    
        Username string 
        Password string 
        Name     string
        MaxIdle  int
    }
    Server struct {
        Host string 
        Port int    
    }
    WebsiteUrl string
}

func main() {
    // Example environment variables
    // DATABASE_HOST=localhost
    // DATABASE_PORT=3306
    // DATABASE_USERNAME=root
    // DATABASE_PASSWORD=secret
    // DATABASE_NAME=example
    // DATABASE_MAX_IDLE=10
    // SERVER_HOST=localhost
    // SERVER_PORT=8080
    // WEBSITE_URL=http://localhost:8080

    // Following lines will load environment variables into Config struct
    // Field delimiter is underscore(_) by default

    var config Config
    envLoader := envloader.New()

    err := envLoader.Load(&config)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", config)
}
```

## Custom Field Names
```go
package main

import (
    "fmt"
    "github.com/metinorak/envloader"
)

// An example nested struct
type Config struct {
    Database struct {
        Host     string     `env:"Host"`
        Port     int        `env:"Port"`
        Username string     `env:"Username"`
        Password string     `env:"Password"`
        Name     string     `env:"Name"`
        MaxIdle  int        `env:"MaxIdle"`
    }  `env:"Database"`
    Server struct {
        Host string         `env:"Host"`
        Port int            `env:"Port"`
    } `env:"Server"`
    WebsiteUrl string `env:"Website"`
}

func main() {
    // Example environment variables
    // Database_Host=localhost
    // Database_Port=3306
    // Database_Username=root
    // Database_Password=secret
    // Database_Name=example
    // Database_MaxIdle=10
    // Server_Host=localhost
    // Server_Port=8080
    // Website=http://localhost:8080

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
