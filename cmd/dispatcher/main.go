package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
)

func main() {
	var oopspath = flag.String("oopspath", "/tmp", "Path to store oopses.")
	var raddr = flag.String("raddr", "127.0.0.1:5678", "IP address of the receiver.")
	flag.Parse()

	quit := make(chan struct{})
	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	fs, err := ioutil.ReadDir(*oopspath)
	if nil != err {
		log.Fatalf(err.Error())
	}
	if len(fs) == 0 {
		log.Println("no oops file found")
		os.Exit(0)
	}

	select {
	case msg := <-term:
		log.Println(msg)
		close(quit)
		os.Exit(0)
	default:
		// FIXME send each oops file
		conn, err := net.Dial("tcp", *raddr)
		if err != nil {
			log.Fatalf(err.Error())
		}
		for _, of := range fs {
			f, err := os.Open(path.Join(*oopspath, of.Name()))
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(f, conn)
			if err != nil {
				log.Printf("error send oops: %s", err)
			}

		}
		conn.Close()

	}
}
