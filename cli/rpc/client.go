package rpc

import "net/rpc"

var serverAddress = ""

type Client struct {
	client *rpc.Client
}

func NewRpcClient(uri string) (*Client, error) {
	if uri == "" {
		uri = "127.0.0.1"
	}

	client, err := rpc.DialHTTP("tcp", serverAddress+":1234")

	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

// args := &server.Args{7,8}
// var reply int
// err = client.Call("Arith.Multiply", args, &reply)
// if err != nil {
// 	log.Fatal("arith error:", err)
// }
// fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

// // Asynchronous call
// quotient := new(Quotient)
// divCall := client.Go("Arith.Divide", args, quotient, nil)
// replyCall := <-divCall.Done	// will be equal to divCall
// check errors, print, etc.
