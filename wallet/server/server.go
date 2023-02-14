package server

import (
	"encoding/hex"
	"fmt"
	"goXdagWallet/xdago/cryptography"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

var HelloServiceName = "Xdag"

type Xdag struct{}

func (s *Xdag) UnlockWallet(request string, replay *string) error {
	b, _ := hex.DecodeString(request)
	hash := cryptography.HashTwice(b)
	*replay = hex.EncodeToString(hash[:])
	return nil
}

func (s *Xdag) LockWallet(request string, replay *string) error {
	b, _ := hex.DecodeString(request)
	hash := cryptography.HashTwice(b)
	*replay = hex.EncodeToString(hash[:])
	return nil
}

func (s *Xdag) Transfer(request string, replay *string) error {
	b, _ := hex.DecodeString(request)
	hash := cryptography.HashTwice(b)
	*replay = hex.EncodeToString(hash[:])
	return nil
}

func (s *Xdag) GetBalance(request string, replay *string) error {
	fmt.Println(request)
	*replay = "1024.000000000"
	return nil
}

func RunServer(ip, port string) {
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

	http.ListenAndServe(ip+":"+port, nil)
}
