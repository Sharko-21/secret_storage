package main

import (
	"main/server"
)

func main() {
	s := server.NewServer()
	<-s.Start()
}
