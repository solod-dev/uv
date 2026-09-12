// Using the libuv library to run a TCP echo server on port 7010.
//
// Usage:
//
//	make example name=echo
//	./build/echo
//	echo 'hello' | nc localhost 7010
//
// Source: https://github.com/libuv/libuv/blob/v1.x/docs/code/tcp-echo-server/main.c
package main

import (
	"unsafe"

	"solod.dev/so/c"
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/os"
	"solod.dev/uv/libuv"
)

const (
	DefaultPort    = 7010
	DefaultBacklog = 128
)

var loop *libuv.Loop_t
var addr libuv.Sockaddr_in

type WriteReq struct {
	req libuv.Write_t
	buf libuv.Buf_t
}

func (wr *WriteReq) Free() {
	buf := unsafe.Slice(wr.buf.Base, wr.buf.Len)
	mem.FreeSlice(mem.System, buf)
	mem.Free(mem.System, wr)
}

func allocBuffer(handle *libuv.Handle_t, suggestedSize c.Size, buf *libuv.Buf_t) {
	_ = handle
	blen := int(suggestedSize)
	b := mem.AllocSlice[c.Char](mem.System, blen, blen)
	buf.Base = c.SliceData[c.Char](b)
	buf.Len = suggestedSize
}

func onClose(handle *libuv.Handle_t) {
	mem.Free(mem.System, handle)
}

func echoWrite(req *libuv.Write_t, status c.Int) {
	if status != 0 {
		msg := c.String(libuv.Strerror(status))
		fmt.Fprintf(os.Stderr, "Write error %s\n", msg)
	}
	wr := c.PtrAs[WriteReq](req)
	wr.Free()
}

func echoRead(client *libuv.Stream_t, nread c.SSize, buf *libuv.ConstBuf_t) {
	if nread > 0 {
		wr := mem.Alloc[WriteReq](mem.System)
		wr.buf = libuv.Buf_init(buf.Base, c.UInt(nread))
		libuv.Write(&wr.req, client, &wr.buf, 1, echoWrite)
		return
	}
	if nread < 0 {
		if nread != libuv.EOF {
			msg := c.String(libuv.Err_name(c.Int(nread)))
			fmt.Fprintf(os.Stderr, "Read error %s\n", msg)
		}
		libuv.Close(c.PtrAs[libuv.Handle_t](client), onClose)
	}

	b := unsafe.Slice(buf.Base, buf.Len)
	mem.FreeSlice(mem.System, b)
}

func onNewConnection(server *libuv.Stream_t, status c.Int) {
	if status < 0 {
		msg := c.String(libuv.Strerror(status))
		fmt.Fprintf(os.Stderr, "New connection error %s\n", msg)
		return
	}

	client := mem.Alloc[libuv.Tcp_t](mem.System)
	libuv.Tcp_init(loop, client)
	if libuv.Accept(server, c.PtrAs[libuv.Stream_t](client)) == 0 {
		libuv.Read_start(c.PtrAs[libuv.Stream_t](client), allocBuffer, echoRead)
	} else {
		libuv.Close(c.PtrAs[libuv.Handle_t](client), onClose)
	}
}

func main() {
	loop = libuv.Default_loop()

	var server libuv.Tcp_t
	libuv.Tcp_init(loop, &server)

	libuv.Ip4_addr("0.0.0.0", DefaultPort, &addr)
	libuv.Tcp_bind(&server, c.PtrAs[libuv.Sockaddr](addr), 0)
	r := libuv.Listen(c.PtrAs[libuv.Stream_t](&server), DefaultBacklog, onNewConnection)
	if r != 0 {
		msg := c.String(libuv.Strerror(r))
		fmt.Fprintf(os.Stderr, "Listen error %s\n", msg)
		os.Exit(1)
	}
	libuv.Run(loop, libuv.RUN_DEFAULT)
}
