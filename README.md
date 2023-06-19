# EnvLoader
## Description
EnvLoader is a simple library that allows you to load environment variables into a model struct. And it supports nested structs too!

## Features
- Supports nested structs
- Field names are converted to upper snake case by default
- But you can specify custom field names with `env` tags
- Field delimiter is dot by default, but you can change it with WithEnvFieldDelimiter option while creating EnvLoader instance

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
    // DATABASE.HOST=localhost
    // DATABASE.PORT=3306
    // DATABASE.USERNAME=root
    // DATABASE.PASSWORD=secret
    // DATABASE.NAME=example
    // DATABASE.MAX_IDLE=10
    // SERVER.HOST=localhost
    // SERVER.PORT=8080
    // WEBSITE_URL=http://localhost:8080

    // Following lines will load environment variables into Config struct
    // Field delimiter is dot(.) by default

    var config Config
    envLoader := envloader.New()

    err := envLoader.Load(&config)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", config)
}
```

## With Options
```go
package main

import (
    "fmt"
    "github.com/metinorak/envloader"
)

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
    // Database*Host=localhost
    // Database*Port=3306
    // Database*Username=root
    // Database*Password=secret
    // Database*Name=example
    // Database*MaxIdle=10
    // Server*Host=localhost
    // Server*Port=8080
    // Website=http://localhost:8080

    // Following lines will load environment variables into Config struct

    var config Config
    envLoader := envloader.New(WithEnvFieldDelimiter("*"))

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
