// Using the libuv library to run an event loop.
//
// Usage:
//
//	make example name=loop
//	./build/loop
//
// Source: https://github.com/libuv/libuv/blob/v1.x/docs/code/default-loop/main.c
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
