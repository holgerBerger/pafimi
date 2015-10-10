package main

import (
	"fmt"
	"github.com/holgerBerger/pafimi/config"
	"log"
	"net/rpc"
	"os"
)

func main() {

	// read config
	config.ReadConf()

	// connect to rpc server
	client, err := rpc.Dial("tcp", config.Conf.Client.Servers[0])
	if err != nil {
		log.Fatal("could not connect to server.")
	}

	// do RPC call
	var reply string
	request := config.Request{User: "bla", Src: os.Args[1], Dst: os.Args[2]}
	err = client.Call("PafimiServerT.AddRequest", request, &reply)
	if err != nil {
		log.Fatal(err)
	}
	// print jobid or error message
	fmt.Println(reply)
	_ = err
}
