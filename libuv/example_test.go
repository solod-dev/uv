package libuv_test

import (
	"solod.dev/uv/libuv"
)

func ExampleRun() {
	loop := libuv.Default_loop()
	libuv.Run(loop, libuv.RUN_DEFAULT)
	libuv.Loop_close(loop)
}
