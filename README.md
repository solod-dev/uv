# uv

Solod bindings for [libuv](https://libuv.org/), a multi-platform library for asynchronous I/O.

## Usage

1. Install libuv for your operating system.

2. Install the Solod bindings.

```
go get solod.dev/uv@latest
```

3. Use it in your code.

```go
package main

import (
	"solod.dev/so/fmt"
	"solod.dev/uv/libuv"
)

func main() {
	loop := libuv.Default_loop()
	fmt.Println("Default loop")
	libuv.Run(loop, libuv.RUN_DEFAULT)
	libuv.Loop_close(loop)
}
```

## Examples

[Event loop](example/loop/main.go)

[Echo server](example/echo/main.go)
