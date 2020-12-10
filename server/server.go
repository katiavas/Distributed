

package main

import (
"flag"
"fmt"
"log"
"net"
"net/rpc"
"uk.ac.bris.cs/gameoflife/gol"
)

func main(){
	pAddr := flag.String("port","8030","Port to listen on")
	flag.Parse()
	err :=rpc.Register(&gol.EngineOperations{})
	if err!=nil{
		log.Fatal("Format of EngineOperations isnt correct",err)
	}
	listener, e := net.Listen("tcp", ":"+*pAddr)
	if e!=nil{
		log.Fatal("Listen error",e)
	}
	fmt.Print("Server is running...")
	defer listener.Close()
	rpc.Accept(listener)
}

