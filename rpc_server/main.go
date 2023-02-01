package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"unsafe"
)

const (
	bits2mime = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var HelloServiceName = "Xdag"

type Xdag struct{}

func (s *Xdag) SendRawTransaction(request string, replay *string) error {
	fmt.Println(request)
	*replay = blockHash(request)
	return nil
}

func (s *Xdag) GetBalance(request string, replay *string) error {
	fmt.Println(request)
	*replay = "1024.000000000"
	return nil
}

func main() {
	rpc.RegisterName(HelloServiceName, new(Xdag))

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     rw,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn)) //rpc.ServeRequest类似于rpc.ServeCodec,以同步的方式处理请求
	})

	http.ListenAndServe("127.0.0.1:10001", nil)
}

func blockHash(block string) string {
	b, _ := hex.DecodeString(block)
	hash := sha256.Sum256(b)
	hash = sha256.Sum256(hash[:])
	return Hash2Address(hash)
}

func Hash2Address(h [32]byte) string {
	address := make([]byte, 32)
	var c, d, j uint
	// every 3 bytes(24 bits) hashs convert to 4 chars(6 bit each)
	// first 24 bytes hash to 32 byte address, ignore last 8 bytes of hash
	for i := 0; i < 32; i++ {
		if d < 6 {
			d += 8
			c <<= 8
			c |= uint(h[j])
			j++
		}
		d -= 6
		address[i] = bits2mime[c>>d&0x3F]
	}
	return bytes2str(address)
}

// unsafe and fast convert bytes slice to string
func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
