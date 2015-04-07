package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/vds/oopsio/receiver"
)

var (
	ip = "localhost:5678"
)

func main() {

	var path = flag.String("path", "/tmp", "Path to store oopses.")
	flag.Parse()

	quit := make(chan struct{})
	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	laddr, err := net.ResolveTCPAddr("tcp", ip)
	if nil != err {
		log.Fatalln(err)
	}
	receiver := receiver.NewReceiver(quit, *path)
	go receiver.Run(laddr)

	select {
	case msg := <-term:
		log.Println(msg)
		close(quit)
	}

	err = receiver.Stop()
	os.Exit(0)
}
