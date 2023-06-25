# EnvLoader
## Description
EnvLoader is a simple library that allows you to load environment variables into a model struct. And it supports nested structs too!

## Features
- Supports nested structs
- Field names are converted to upper snake case by default
- Custom field names can be defined with `env` tag
- Default values can be defined with `default` tag
- Supports slice and map types
- Requirement check can be enabled with `required` tag like `required:"true"`. It is disabled by default.
- Nested struct fields' variable names consist of parent struct name and field name. For example, `Database_Host` for `Database struct { Host string }`
- Field delimiter is underscore(_) by default. It can be disabled using ``env:"-"``. In this case struct field names will not contain parent struct name. For example, `Host` for `Database struct { Host string }`

## Installation
```bash
go get github.com/metinorak/envloader
```

## Basic Usage
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
    FormulaConstants map[string]float64
    Proxies []string
}

func main() {
    // Example environment variables for the above struct
    // DATABASE_HOST=localhost
    // DATABASE_PORT=3306
    // DATABASE_USERNAME=root
    // DATABASE_PASSWORD=secret
    // DATABASE_NAME=example
    // DATABASE_MAX_IDLE=10
    // SERVER_HOST=localhost
    // SERVER_PORT=8080
    // WEBSITE_URL=http://localhost:8080
    // FORMULA_CONSTANTS=pi:3.14,e:2.71
    // PROXIES=example.com:8080,example2.com:8080

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

## With Tags
```go
package main

import (
    "fmt"
    "github.com/metinorak/envloader"
)

// An example nested struct
type Config struct {
    Database struct {
        Host     string     `env:"Host" default:"localhost"`
        Port     int        `env:"Port" default:"5000"`
        Username string     `env:"Username" required:"true"`
        Password string     `env:"Password" required:"true"`
        Name     string     `env:"Name" default:"example" required:"true"`
        MaxIdle  int        `env:"MaxIdle"`
    }  `env:"Database"`
    Server struct {
        Host string         `env:"Host"`
        Port int            `env:"Port"`
    } `env:"Server"`
    WebsiteUrl string       `env:"Website"`
    FormulaConstant float64 `env:"FormulaConstant" default:"3.14"`
}

func main() {
    // Example environment variables for the above struct
    // Database_Host=localhost
    // Database_Port=3306
    // Database_Username=root
    // Database_Password=secret
    // Database_Name=example
    // Database_MaxIdle=10
    // Server_Host=localhost
    // Server_Port=8080
    // Website=http://localhost:8080
    // FormulaConstant=3.14

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

## With Disabled Parent Struct Name and Fields
```go
package main

import (
    "fmt"
    "github.com/metinorak/envloader"
)

// An example nested struct
type Config struct {
    Database struct {
        Host     string     `env:"-" default:"localhost"`
        Port     int        `env:"-" default:"5000"`
        Username string     `env:"Username" required:"true"`
        Password string     `env:"Password" required:"true"`
        Name     string     `env:"Name" default:"example" required:"true"`
        MaxIdle  int        `env:"MaxIdle"`
    }  `env:"-"`
    Server struct {
        Host string         `env:"Host"`
        Port int            `env:"Port"`
    } `env:"Server"`
    WebsiteUrl string       `env:"Website"`
    FormulaConstant float64 `env:"-"`
}

func main() {
    // Example environment variables for the above struct
    // Username=root
    // Password=secret
    // Name=example
    // MaxIdle=10
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
}
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
