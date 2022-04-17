# IniGo

Read INI files in Go

(writing coming soon)

## Installation

If you're using Go modules, just import IniGo

```golang
import "github.com/anthony-y/inigo"
```

and it will be built when you run

```
go build
```

If your code is in $GOPATH/src, then you'll need to first run

```
go get "github.com/anthony-y/inigo"
```

## Getting started

Given this INI file

```ini
[MyData]
myVariable="Hello world"
```

The following code will output

```
Hello world
```

```golang
package main

import (
    "fmt"

    "github.com/anthony-y/inigo"
)

func main() {
    ini, errs := inigo.ReadIni("example.ini")
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
		return
    }
    
    fmt.Println(ini["MyData"]["myVariable"])
}

```